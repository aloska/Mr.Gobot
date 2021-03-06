package agent

import (
	"github.com/gin-gonic/gin"
	"github.com/pterm/pterm"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/ini.v1"
	"io/ioutil"
	"net/url"
	"os"
	"regexp"
	"runtime"
	"time"
)

type Agent struct{
	path string		//путь к организму
	o *Organism		//собсно организм
	server *gin.Engine //сервер http
	port int    //порт сервера
	config *ini.File //ини-конфигурация
	logfile string //путь к лог-файлу (лежит в папке организма /vm/logs/)
	log *zap.Logger

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
	a.o=&Organism{}
	if! a.o.Init(a){
		a.fatalout()
		return
	}



	//все прошло удачно

}

func (a* Agent) checkpaths() bool{
	pterm.DefaultSection.Println("Проверка структуры папок и файлов...")

	pslice:=[]string{"/Senses", "/Actions", "/Brain", "/Vegetatic"}

	p, _ := pterm.DefaultProgressbar.WithTotal(len(pslice)).WithTitle("Проверяем папки").Start()
	for i := 0; i < p.Total; i++ {
		p.Title = "Папка " + pslice[i]
		if  _, err :=  os.Stat(a.path+pslice[i]); os.IsNotExist(err) {
			pterm.Error.Println("Нет папки "+pslice[i])
			return false
		} else {
			pterm.Success.Println("Есть папка " + pslice[i])
			p.Increment()
		}
		time.Sleep(time.Second / 2)
	}
	p.Stop()


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
					serr:="В /Actions/Effector-xxx/ не может быть "+file.Name()
					pterm.Warning.Println("Что это за файл? "+file.Name())
					pterm.Error.Println(serr)
					a.log.Error(serr,zap.String("/Vegetatic/"+inp, "Плохое имя файла "+file.Name()))
					return false
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
	pterm.Info.Println("Похоже, все папки и файлы Организма в порядке... Но это не точно! ")
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
	pterm.Info.Println("Концигурационный файл в порядке!")
	return  true
}
