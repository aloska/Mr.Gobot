package agent

import (
	"context"
	"fmt"
	"github.com/cheggaaa/pb/v3"
	"github.com/gin-gonic/gin"
	"github.com/pterm/pterm"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/ini.v1"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"runtime"
	"strconv"
	"sync"
	"time"
)


type Agent struct{
	path string		//путь к организму
	o *Organism		//собсно организм
	router *gin.Engine //роутер http
	server *http.Server //сервер http
	port int    //порт сервера
	config *ini.File //ини-конфигурация
	logfile string //путь к лог-файлу (лежит в папке организма /vm/logs/)
	log *zap.Logger

	doatstart string //запись в конйигурационном файле, как начать жить "pause", "live"
	//каналы для управления организмом
	live chan struct{}	//комманда на жизнь
	pause chan struct{} //пауза
	quit chan struct{}  //выключаемся
	wga sync.WaitGroup	//wait-group для гороутин, запущенных агентом
}

func (a* Agent) Live (pathtoOrganism string){
	a.path=pathtoOrganism
	//издрасти
	a.welcome()
	//читаем конфиг
	if !a.readConfig() {
		a.fatalout()
		return
	}
	//создаем логи
	if !a.startLog(){
		a.fatalout()
		return
	}
	defer a.log.Sync()
	a.log.Info("Старт организма", zap.String("path", a.path))

	//проверка папок и файлов
	if !a.checkpaths(){
		a.fatalout()
		return
	}

	//Инициализируем организм
	pterm.DefaultSection.Println("Инициализация...")
	a.o=&Organism{}
	if! a.o.Init(a){
		a.fatalout()
		return
	}
	pterm.Success.Println("Инициализация прошла успешно")

	//Проверка работоспособности
	pterm.DefaultSection.Println("Проверка работоспособности...")
	tmpl := `{{ yellow "Пока все хорошо:" }} {{ bar . "[" "=" (cycle . "↖" "↗" "↘" "↙" ) "." "]"}} {{speed . | rndcolor }} {{percent .}}  {{string . "Проверка..." | green}}`
	// start bar based on our template
	bar := pb.ProgressBarTemplate(tmpl).Start(a.o.countall)
	//канал, для получения данных от организма о состоянии проверки
	c := make(chan error)
	//погнали
	go a.o.Check(a,c)
	//читаем из канала до закрытия
	isok:=true
	for err:= range c {
		bar.Increment()
		if err!=nil{
			a.log.Error(err.Error())
			isok=false
		}
	}
	bar.Finish()
	println()
	if !isok{
		a.errorr("Во время работы обнаружены ошибки. Смотри лог.")
		a.fatalout()
		return
	}
	time.Sleep(time.Second)
	fmt.Println()

	//все прошло удачно
	pterm.DefaultSection.Println("Go Live!")

	//каналы для контроля организмом
	a.live = make(chan struct{})
	a.quit = make(chan struct{})
	a.pause = make(chan struct{})

	//запускаем асинхронно организм
	a.wga.Add(1)
	go a.o.Live()
	switch a.doatstart {
	case "pause":
		a.pause <- struct{}{}
		break
	default:
		a.live <- struct{}{}
	}

	//запускаем http-server
	pterm.DefaultSection.Println("Запускаем http-server")
	a.router=gin.Default()
	a.makeRoutes()
	//эти пляски для нормального выключения сервера по srv.shutdown
	a.server = &http.Server{
		Addr:    ":"+strconv.Itoa(a.port),
		Handler: a.router,
	}
	//здесь мы блокируемся
	if err := a.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		a.log.Fatal(fmt.Sprintf( "listen: %s\n", err))
	}

	//а здесь конец работы
	a.welcome()
	pterm.Println("Ну пока!")

}

func (a* Agent) makeRoutes(){
	a.router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Welcome Okagamga2.0 agent")
	})
	a.router.GET("/quitall", a.routeQuitall)
	a.router.GET("/step", a.routeStep)
}

func (a *Agent) routeStep (c *gin.Context){
	if a.o.state!="pause" {
		a.pause <- struct{}{}
	}
	time.Sleep(time.Second/10)
	strt:=time.Now()
	a.o.Step()
	elapsed := time.Since(strt)
	c.String(http.StatusOK, fmt.Sprintf("Step complete in %v", elapsed))
}

func (a *Agent) routeQuitall(c* gin.Context){
	//сначала сохраним организм
	if a.o.state!="pause" {
		a.pause <- struct{}{}
	}
	time.Sleep(time.Second)
	a.o.Sleep()
	a.quit<- struct{}{}
	a.wga.Wait() //подождем, пока организм не выключится

	close(a.live)
	close(a.pause)
	close(a.quit)

	c.String(http.StatusOK, "Bye Okagamga2.0 agent")
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	//ctx - контекст нужен для того, чтобы сервер успел отправить что не отправил
	if err := a.server.Shutdown(ctx); err != nil {
		a.log.Warn(fmt.Sprintf( "listen: %s\n", err))
	}
}

func (a* Agent) checkpaths() bool{
	pterm.DefaultSection.Println("Проверка структуры папок и файлов...")

	synverify:=map[int]int{}//для проверки уникальности номеров синаптических полей

	pslice:=[]string{"/Senses", "/Actions", "/Brain", "/Vegetatic"}
	for i := 0; i < 4; i++ {
		if  _, err :=  os.Stat(a.path+pslice[i]); os.IsNotExist(err) {
			pterm.Error.Println("Нет папки "+pslice[i])
			return false
		} else {
			pterm.Success.Println("Есть папка " + pslice[i])
		}
	}

	pterm.DefaultSection.Println("Проверка папки /Senses...")
	inputs:=[]string{}
	files, _ := ioutil.ReadDir(a.path+"/Senses")
	for _, file := range files {
		if file.IsDir(){
			//должны быть только папки Input-xxx
			if match, _ := regexp.MatchString("Input-[0-9]+",file.Name()); !match{
				serr:="Все папки в директории /Senses должны быть формата /Input-X, где X - число"
				pterm.Error.Println(serr)
				a.log.Error(serr,zap.String("/Senses", "Плохое имя папки "+file.Name()))
				return false
			}
			inputs = append(inputs,file.Name())
		} else{
			//должны быть только файлы  Neuron-XXX.neurons GenNeuron-XXX.genes Synapse-XXX.chemical...
			if match, _ := regexp.MatchString(
				"(syn-[0-9]+x[0-9]+.[0-9]+|Neuron-[0-9]+.neurons|Synapse-[0-9]+.chemical|GenNeuron-[0-9]+.genes)",
					file.Name()); !match{
						serr:="В корне /Senses могут быть файлы формата syn-NxM.c Neuron-XXX.neurons GenNeuron-XXX.genes Synapse-XXX.chemical"
						pterm.Warning.Println("Что это за файл? "+file.Name())
						pterm.Error.Println(serr)
						a.log.Error(serr,zap.String("/Senses", "Плохое имя файла "+file.Name()))
						return false
			}
			//если это синаптическое поле - занесем его номер и проверим, не было ли такого уже
			if match, _:= regexp.MatchString("syn-[0-9]+x[0-9]+.[0-9]+",file.Name()); match{
				//отпарсим номер и размер
				re := regexp.MustCompile(`[0-9]+`)
				ss:=re.FindAllString(file.Name(), -1)
				num,_:=strconv.Atoi(ss[2])
				a.info("Обнаружено синаптическое поле "+file.Name()+" в /Senses")
				if _, found := synverify[num]; found{
					serr:="Номера синаптических полей совпадают! "+file.Name()
					pterm.Error.Println(serr)
					a.log.Error(serr,zap.String("/Senses", "Плохое имя файла "+file.Name()))
					return false
				}
				synverify[num]=1
			}

		}
	}
	pterm.Info.Println("/Senses выглядит сносно")

	pterm.DefaultSection.Println("Проверка папок /Senses/Input-...")
	for _, inp:=range inputs{
		gendatagenes:=false
		genreceptorgenes:=false
		files, _ := ioutil.ReadDir(a.path+"/Senses/"+inp)
		for _, file := range files {
			if file.IsDir(){
				//здесь не должно быть папок
				serr:="В /Senses/Input-xxx/ не должно быть папок"
				pterm.Warning.Println("Что это за папка? "+file.Name())
				pterm.Error.Println(serr)
				a.log.Error(serr,zap.String("/Senses/"+inp, "не известная папка"))
				return false

			} else{
				//должны быть только файлы
				pat:="(syn-[0-9]+x[0-9]+.[0-9]+|Data.data|GenData.genes|GenReceptor-[0-9]+.genes|Receptor-[0-9]+.receptors|Neuron-[0-9]+.neurons|Synapse-[0-9]+.chemical|GenNeuron-[0-9]+.genes)"
				if match, _ := regexp.MatchString(pat, file.Name()); !match{
					serr:="В /Senses/Input-xxx/ не может быть "+file.Name()
					pterm.Warning.Println("Что это за файл? "+file.Name())
					pterm.Error.Println(serr)
					a.log.Error(serr,zap.String("/Senses/"+inp, "Плохое имя файла "+file.Name()))
					return false
				}
				//если это синаптическое поле - занесем его номер и проверим, не было ли такого уже
				if match, _:= regexp.MatchString("syn-[0-9]+x[0-9]+.[0-9]+",file.Name()); match{
					//отпарсим номер и размер
					re := regexp.MustCompile(`[0-9]+`)
					ss:=re.FindAllString(file.Name(), -1)
					num,_:=strconv.Atoi(ss[2])
					a.info("Обнаружено синаптическое поле "+file.Name()+" в /Senses/"+inp)
					if _, found := synverify[num]; found{
						serr:="Номера синаптических полей совпадают! "+file.Name()
						pterm.Error.Println(serr)
						a.log.Error(serr,zap.String("/Senses/"+inp, "Плохое имя файла "+file.Name()))
						return false
					}
					synverify[num]=1
				}

				//GenData.genes должен быть обязательно
				if match, _ := regexp.MatchString("GenData.genes", file.Name()); match{
					gendatagenes=true
				}
				//Какой-нибудь GenReceptor-xxx.genes должен быть обязательно
				if match, _ := regexp.MatchString("GenReceptor-[0-9]+.genes", file.Name()); match{
					genreceptorgenes=true
				}
			}
		}
		if gendatagenes && genreceptorgenes{
			pterm.Info.Println("/Senses/"+inp+" выглядит сносно")
		} else{
			serr:="В /Senses/Input-xxx/ должен быть GenReceptor-0.genes и GenData.genes обязательно"
			pterm.Error.Println(serr)
			a.log.Error(serr,zap.String("/Senses/"+inp, "Не хватает файла "))
			return false
		}
	}

	pterm.DefaultSection.Println("Проверка папки /Actions...")
	inputs=[]string{}
	files, _ = ioutil.ReadDir(a.path+"/Actions")
	for _, file := range files {
		if file.IsDir(){
			//должны быть только папки Effector-xxx
			if match, _ := regexp.MatchString("Effector-[0-9]+",file.Name()); !match{
				serr:="Все папки в директории /Actions должны быть формата /Effector-X, где X - число"
				pterm.Error.Println(serr)
				a.log.Error(serr,zap.String("/Actions", "Плохое имя папки "+file.Name()))
				return false
			}
			inputs = append(inputs,file.Name())
		} else{
			//должны быть только файлы  Neuron-XXX.neurons GenNeuron-XXX.genes Synapse-XXX.chemical...
			if match, _ := regexp.MatchString(
				"(syn-[0-9]+x[0-9]+.[0-9]+|Neuron-[0-9]+.neurons|Synapse-[0-9]+.chemical|GenNeuron-[0-9]+.genes)",
				file.Name()); !match{
				serr:="В корне /Actions могут быть файлы формата syn-NxM.c Neuron-XXX.neurons GenNeuron-XXX.genes Synapse-XXX.chemical"
				pterm.Warning.Println("Что это за файл? "+file.Name())
				pterm.Error.Println(serr)
				a.log.Error(serr,zap.String("/Actions", "Плохое имя файла "+file.Name()))
				return false
			}
			//если это синаптическое поле - занесем его номер и проверим, не было ли такого уже
			if match, _:= regexp.MatchString("syn-[0-9]+x[0-9]+.[0-9]+",file.Name()); match{
				//отпарсим номер и размер
				re := regexp.MustCompile(`[0-9]+`)
				ss:=re.FindAllString(file.Name(), -1)
				num,_:=strconv.Atoi(ss[2])
				a.info("Обнаружено синаптическое поле "+file.Name()+" в /Actions")
				if _, found := synverify[num]; found{
					serr:="Номера синаптических полей совпадают! "+file.Name()
					pterm.Error.Println(serr)
					a.log.Error(serr,zap.String("/Actions", "Плохое имя файла "+file.Name()))
					return false
				}
				synverify[num]=1
			}
		}
	}
	pterm.Info.Println("/Actions выглядит сносно")

	pterm.DefaultSection.Println("Проверка папок /Actions/Effector-...")
	for _, inp:=range inputs{
		gen1:=false
		gen2:=false
		files, _ := ioutil.ReadDir(a.path+"/Actions/"+inp)
		for _, file := range files {
			if file.IsDir(){
				//здесь не должно быть папок
				serr:="В /Actions/Effector-xxx/ не должно быть папок"
				pterm.Warning.Println("Что это за папка? "+file.Name())
				pterm.Error.Println(serr)
				a.log.Error(serr,zap.String("/Actions/"+inp, "не известная папка"))
				return false

			} else{
				//должны быть только файлы  Neuron-XXX.neurons GenNeuron-XXX.genes Synapse-XXX.chemical
				pat:="(syn-[0-9]+x[0-9]+.[0-9]+|DataOut.data|GenDataOut.genes|GenPreffector-[0-9]+.genes|Preffector-[0-9]+.preffectors|Neuron-[0-9]+.neurons|Synapse-[0-9]+.chemical|GenNeuron-[0-9]+.genes)"
				if match, _ := regexp.MatchString(pat, file.Name()); !match{
					serr:="В /Actions/Effector-xxx/ не может быть "+file.Name()
					pterm.Warning.Println("Что это за файл? "+file.Name())
					pterm.Error.Println(serr)
					a.log.Error(serr,zap.String("/Actions/"+inp, "Плохое имя файла "+file.Name()))
					return false
				}
				//если это синаптическое поле - занесем его номер и проверим, не было ли такого уже
				if match, _:= regexp.MatchString("syn-[0-9]+x[0-9]+.[0-9]+",file.Name()); match{
					//отпарсим номер и размер
					re := regexp.MustCompile(`[0-9]+`)
					ss:=re.FindAllString(file.Name(), -1)
					num,_:=strconv.Atoi(ss[2])
					a.info("Обнаружено синаптическое поле "+file.Name()+" в /Actions/"+inp)
					if _, found := synverify[num]; found{
						serr:="Номера синаптических полей совпадают! "+file.Name()
						pterm.Error.Println(serr)
						a.log.Error(serr,zap.String("/Actions/"+inp, "Плохое имя файла "+file.Name()))
						return false
					}
					synverify[num]=1
				}

				//GenDataOut.genes должен быть обязательно
				if match, _ := regexp.MatchString("GenDataOut.genes", file.Name()); match{
					gen1=true
				}
				//GenPreffector-xxx.genes должен быть обязательно
				if match, _ := regexp.MatchString("GenPreffector-[0-9]+.genes", file.Name()); match{
					gen2=true
				}
			}
		}
		if gen1 && gen2{
			pterm.Info.Println("/Actions/"+inp+" выглядит сносно")
		} else{
			serr:="В /Actions/Effector-xxx/ должен быть GenPreffector-xxx.genes и GenDataOut.genes обязательно"
			pterm.Error.Println(serr)
			a.log.Error(serr,zap.String("/Actions/"+inp, "Не хватает файла "))
			return false
		}
	}

	pterm.DefaultSection.Println("Проверка папки /Brain...")
	inputs=[]string{}
	files, _ = ioutil.ReadDir(a.path+"/Brain")
	for _, file := range files {
		if file.IsDir(){
			//должны быть только папки Core-xxx
			if match, _ := regexp.MatchString("Core-[0-9]+",file.Name()); !match{
				serr:="Все папки в директории /Brain должны быть формата /Core-X, где X - число"
				pterm.Error.Println(serr)
				a.log.Error(serr,zap.String("/Brain", "Плохое имя папки "+file.Name()))
				return false
			}
			inputs = append(inputs,file.Name())
		} else{
			//не должно быть файлов в корне
			serr:="В корне /Brain не должно быть никаких файлов"
			pterm.Warning.Println("Что это за файл? "+file.Name())
			pterm.Error.Println(serr)
			a.log.Error(serr,zap.String("/Brain", "Не должно быть файла "+file.Name()))
			return false
		}
	}
	pterm.Info.Println("/Brain выглядит сносно")

	pterm.DefaultSection.Println("Проверка папок /Brain/Core-...")
	for _, inp:=range inputs{
		gen1, gen2:=false,false
		files, _ := ioutil.ReadDir(a.path+"/Brain/"+inp)
		for _, file := range files {
			if file.IsDir(){
				//здесь не должно быть папок
				serr:="В /Brain/Core-xxx/ не должно быть папок"
				pterm.Warning.Println("Что это за папка? "+file.Name())
				pterm.Error.Println(serr)
				a.log.Error(serr,zap.String("/Brain/"+inp, "не известная папка"))
				return false

			} else{
				//должны быть только файлы  Neuron-XXX.neurons GenNeuron-XXX.genes Synapse-XXX.chemical
				pat:="(syn-[0-9]+x[0-9]+.[0-9]+|Neuron-[0-9]+.neurons|Synapse-[0-9]+.chemical|GenNeuron-[0-9]+.genes)"
				if match, _ := regexp.MatchString(pat, file.Name()); !match{
					serr:="В /Brain/Core-xxx/ не может быть "+file.Name()
					pterm.Warning.Println("Что это за файл? "+file.Name())
					pterm.Error.Println(serr)
					a.log.Error(serr,zap.String("/Brain/"+inp, "Плохое имя файла "+file.Name()))
					return false
				}
				//если это синаптическое поле - занесем его номер и проверим, не было ли такого уже
				if match, _:= regexp.MatchString("syn-[0-9]+x[0-9]+.[0-9]+",file.Name()); match{
					//отпарсим номер и размер
					re := regexp.MustCompile(`[0-9]+`)
					ss:=re.FindAllString(file.Name(), -1)
					num,_:=strconv.Atoi(ss[2])
					a.info("Обнаружено синаптическое поле "+file.Name()+" в /Brain/"+inp)
					if _, found := synverify[num]; found{
						serr:="Номера синаптических полей совпадают! "+file.Name()
						pterm.Error.Println(serr)
						a.log.Error(serr,zap.String("/Brain/"+inp, "Плохое имя файла "+file.Name()))
						return false
					}
					synverify[num]=1
				}

				//GenNeuron-xxx.genes и syn-[0-9]+x[0-9]+.[0-9]+ должен быть обязательно
				if match, _ := regexp.MatchString("GenNeuron-[0-9]+.genes", file.Name()); match{
					gen1=true
				}
				if match, _ := regexp.MatchString("syn-[0-9]+x[0-9]+.[0-9]+", file.Name()); match{
					gen2=true
				}

			}
		}
		if gen1 && gen2{
			pterm.Info.Println("/Brain/"+inp+" выглядит сносно")
		} else{
			serr:="В /Brain/Core-xxx/ должен быть GenNeuron-xxx.genes  и syn-[0-9]+x[0-9]+.[0-9]+"
			pterm.Error.Println(serr)
			a.log.Error(serr,zap.String("/Brain/"+inp, "Не хватает файла "))
			return false
		}
	}

	pterm.DefaultSection.Println("Проверка папки /Vegetatic...")
	inputs=[]string{}
	files, _ = ioutil.ReadDir(a.path+"/Vegetatic")
	for _, file := range files {
		if file.IsDir(){
			//должны быть только папки Effector-xxx
			if match, _ := regexp.MatchString("Effector-[0-9]+",file.Name()); !match{
				serr:="Все папки в директории /Vegetatic должны быть формата /Effector-X, где X - число"
				pterm.Error.Println(serr)
				a.log.Error(serr,zap.String("/Vegetatic", "Плохое имя папки "+file.Name()))
				return false
			}
			inputs = append(inputs,file.Name())
		} else{
			//должны быть только файлы  Neuron-XXX.neurons GenNeuron-XXX.genes Synapse-XXX.chemical
			pat:="(syn-[0-9]+x[0-9]+.[0-9]+|Neuron-[0-9]+.neurons|Synapse-[0-9]+.chemical|GenNeuron-[0-9]+.genes)"
			if match, _ := regexp.MatchString(pat, file.Name()); !match{
				serr:="В /Vegetatic/ не может быть "+file.Name()
				pterm.Warning.Println("Что это за файл? "+file.Name())
				pterm.Error.Println(serr)
				a.log.Error(serr,zap.String("/Vegetatic/", "Плохое имя файла "+file.Name()))
				return false
			}
			//если это синаптическое поле - занесем его номер и проверим, не было ли такого уже
			if match, _:= regexp.MatchString("syn-[0-9]+x[0-9]+.[0-9]+",file.Name()); match{
				//отпарсим номер и размер
				re := regexp.MustCompile(`[0-9]+`)
				ss:=re.FindAllString(file.Name(), -1)
				num,_:=strconv.Atoi(ss[2])
				a.info("Обнаружено синаптическое поле "+file.Name()+" в /Vegetatic")
				if _, found := synverify[num]; found{
					serr:="Номера синаптических полей совпадают! "+file.Name()
					pterm.Error.Println(serr)
					a.log.Error(serr,zap.String("/Vegetatic", "Плохое имя файла "+file.Name()))
					return false
				}
				synverify[num]=1
			}
		}
	}
	pterm.Info.Println("/Vegetatic выглядит сносно")

	if len(inputs)>0 {
		pterm.DefaultSection.Println("Проверка папок /Vegetatic/Effector-...")
	}
	for _, inp:=range inputs{
		gen1:=false
		gen2:=false
		files, _ := ioutil.ReadDir(a.path+"/Vegetatic/"+inp)
		for _, file := range files {
			if file.IsDir(){
				//здесь не должно быть папок
				serr:="В /Vegetatic/Effector-xxx/ не должно быть папок"
				pterm.Warning.Println("Что это за папка? "+file.Name())
				pterm.Error.Println(serr)
				a.log.Error(serr,zap.String("/Vegetatic/"+inp, "не известная папка"))
				return false

			} else{
				//должны быть только файлы  Neuron-XXX.neurons GenNeuron-XXX.genes Synapse-XXX.chemical
				pat:="(syn-[0-9]+x[0-9]+.[0-9]+|DataOut.data|GenDataOut.genes|GenPreffector-[0-9]+.genes|Preffector-[0-9]+.preffectors|Neuron-[0-9]+.neurons|Synapse-[0-9]+.chemical|GenNeuron-[0-9]+.genes)"
				if match, _ := regexp.MatchString(pat, file.Name()); !match{
					serr:="В /Vegetatic/Effector-xxx/ не может быть "+file.Name()
					pterm.Warning.Println("Что это за файл? "+file.Name())
					pterm.Error.Println(serr)
					a.log.Error(serr,zap.String("/Vegetatic/"+inp, "Плохое имя файла "+file.Name()))
					return false
				}
				//если это синаптическое поле - занесем его номер и проверим, не было ли такого уже
				if match, _:= regexp.MatchString("syn-[0-9]+x[0-9]+.[0-9]+",file.Name()); match{
					//отпарсим номер и размер
					re := regexp.MustCompile(`[0-9]+`)
					ss:=re.FindAllString(file.Name(), -1)
					num,_:=strconv.Atoi(ss[2])
					a.info("Обнаружено синаптическое поле "+file.Name()+" в /Vegetatic/"+inp)
					if _, found := synverify[num]; found{
						serr:="Номера синаптических полей совпадают! "+file.Name()
						pterm.Error.Println(serr)
						a.log.Error(serr,zap.String("/Vegetatic/"+inp, "Плохое имя файла "+file.Name()))
						return false
					}
					synverify[num]=1
				}

				//GenDataOut.genes должен быть обязательно
				if match, _ := regexp.MatchString("GenDataOut.genes", file.Name()); match{
					gen1=true
				}
				//GenPreffector-xxx.genes должен быть обязательно
				if match, _ := regexp.MatchString("GenPreffector-[0-9]+.genes", file.Name()); match{
					gen2=true
				}
			}
		}
		if gen1 && gen2{
			pterm.Info.Println("/Vegetatic/"+inp+" выглядит сносно")
		} else{
			serr:="В /Vegetatic/Effector-xxx/ должен быть GenPreffector-0.genes и GenDataOut.genes обязательно"
			pterm.Error.Println(serr)
			a.log.Error(serr,zap.String("/Vegetatic/"+inp, "Не хватает файла "))
			return false
		}
	}
	pterm.Success.Println("Похоже, все папки и файлы Организма в порядке... Но это не точно! ")
	return true
}

//ибо zap не любит винду TODO - проверить в новой версии
func newWinFileSink(u *url.URL) (zap.Sink, error) {
	// Remove leading slash left by url.Parse()
	return os.OpenFile(u.Path[1:], os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
}

func (a* Agent) startLog() bool {
	//проверяем папки, если их нет, то создаем
	if err:=os.MkdirAll(a.path+"/vm/logs", os.ModePerm); err!=nil{
		pterm.Error.Println("Ошибка создания служебных папок")
		pterm.Error.Println(err)
		return false
	}
	var err error
	var fslice []string

	//ВНИМАНИЕ! Го использует имеенно "2006.01.02 15:04:05" для форматирования. Если поменять эти числа - будет бага!
	if runtime.GOOS == "windows" {
		a.logfile = "winfile:///" + a.path + "/vm/logs/" + time.Now().Format("2006-01-02_15.04.05") + ".log"
	}else{
		a.logfile = a.path + "/vm/logs/" + time.Now().Format("2006-01-02_15.04.05") + ".log"
	}
	//будем логгить в консоль?
	if b, err:=a.config.Section("logger").Key("stdout").Bool(); b && err==nil {
		fslice=append(fslice,"stdout")
	}
	fslice=append(fslice,a.logfile)


	//потому что zap не любит виндовс
	if runtime.GOOS == "windows" {
		zap.RegisterSink("winfile", newWinFileSink)
	}

	cfg:=zap.Config{
		Encoding:    "json",
		Level:       zap.NewAtomicLevelAt(zapcore.DebugLevel),
		OutputPaths: fslice,
		ErrorOutputPaths: []string{"stderr"},
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey: "message",

			LevelKey:    "level",
			EncodeLevel: zapcore.CapitalLevelEncoder,

			TimeKey:    "time",
			EncodeTime: zapcore.ISO8601TimeEncoder,

			CallerKey:    "caller",
			EncodeCaller: zapcore.ShortCallerEncoder,
		},
	}

	a.log, err=cfg.Build()
	if err!=nil{
		a.warning("Ошибка создания логгера")
		a.warning(err.Error())
		return false
	}
	pterm.Info.Println("Логгер запущен!")
	return true
}

func (a* Agent) fatalout(){
	pterm.Fatal.WithFatal(false).Println("Ошибки фатальны! Я выключаюсь...")
}
func (a* Agent) warning(warn string){
	pterm.Warning.Println(warn)
}
func (a* Agent) info(info string){
	pterm.Info.Println(info)
}
func (a* Agent) errorr(er string){
	pterm.Error.Println(er)
}

func (a* Agent) welcome(){
	pterm.DefaultBigText.WithLetters(
		pterm.NewLettersFromStringWithStyle("Okagamga", pterm.NewStyle(pterm.FgCyan)),
		pterm.NewLettersFromStringWithStyle("2.0", pterm.NewStyle(pterm.FgLightMagenta))).
		Render()
}

func (a* Agent) readConfig() bool {
	var err error
	a.config, err = ini.Load(a.path+"/agent.ini")
	if err!=nil{
		pterm.Error.Println("Ошибка конфигурационного файла")
		pterm.Error.Println(err)
		return false
	}
	a.port, err =a.config.Section("server").Key("httpport").Int()
	if err!=nil{
		pterm.Error.Println("Ошибка конфигурационного файла")
		pterm.Error.Println(err)
		return false
	}
	a.doatstart=a.config.Section("organism").Key("doatstart").String()

	pterm.Info.Println("Концигурационный файл в порядке!")
	return  true
}
