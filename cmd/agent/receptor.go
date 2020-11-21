package agent

//Receptor - минимальная единица, передающая вещества от входных данных
/*
Является, по сути переводчиком с языка цифр (информации) на язык веществ

Данные складываются в специальный файл и сразу обрабатываются. Вычисляются их производные по времени и пространству, и вторые производные
и даже скорость.
На каждый из этих параметров (чистых и вычисленных) необходимо поставить по два рецептора - положительный и отрицательный.
Положительный выдает вещество (ацетилхолин?) по своим аксонам пропорцианально величине параметра.
Отрицательный выдает вещество обратно пропорцианально величине параметра.
Например, если данные состоят из байтов:
при значении 200 в данных, рецептор положительный выдаст 200 ацетилхолина, а отрицательный выдаст 255-200=55
при значении 0 в данных, рецептор положительный выдаст 0 ацетилхолина, а отрицательный выдаст 255-0=255

Еще можно разработать рецептор среднего значения, который тем больше плюет, чем ближе к среднему значению.
Да и вообще какой угодно, если угодно, надо экспериментировать))
Таким образом, вещества сигнализируют не только большие данные, а вообще всякие.

У рецепторов аксоны направлены в свой файл синапсов, свойственный данному типу входного сигнала.
Фактически, первичная зрительная кора, например, будет распологаться (и клеточно и синаптически) отдельно от мозга, являясь периферическим входом

Поля сгрупированы по 8 байт
*/
//288 байт
type Receptor struct {
	Threshold uint64 /*порог срабатывания.
	0 (или 0xffffffffffffffff для отрицательного рецептора) означает, что нет порога
	обычно - это величина, ниже (или выше, если рецептор отрицательного типа) которой рецептор молчит
	*/

	Typer ReceptorTypeEnum //Тип рецептора  - четные значения для положительных рецепторов, нечетные - для отрицательных
	SynNumber SynEnum /*номер поля синапсов, в синапсах которого находятся аксоны	 */
	Ndata   uint32 //номер ячейки с данными, на которую нацелен рецептор

	NdataW  byte   //номер бита/байта/слова... из этой ячейки
	NdataWb byte
	Serv1 byte
	Typemedi NeuronTypeEnum //byte тип выплевываемого медиатора (обычно ацетилхолин)
	Force    uint16 //сила реакции
	Divforce uint16 /*обратная сила реакции
	для разных типов означает что-то свое, но чаще всего это коэффициенты количества вещества в ответ на стимул:
	vesiculs = val * force / divforce
	Для рецепторов обратного типа (0xff-val)*f/df
	*/

	Gen uint16/* Номер гена (индекс гена в слайсе genes []GenReceptor), к которому относится рецептор
	Это важный параметр, по которому рецептор может узнать гено-специфические способы работы
	*/
	Serv2 uint16
	//У каждого рецептора до 32 аксонов, которые передают вещества
	A     uint32   //Состояния аксонов. каждый бит отвечает за 1 аксон. Если он 1 - аксон есть и его надо обрабатывать
	Axons [32]Axon //Координаты аксонов в общем файле синапсов, количество везикул и кальций
	//
}

/*DoReceptorUInt32 - запустить обработку ячейки данных. Возвращает количество вещества, которое он вбросил аксонами
 */
func (r *Receptor) DoReceptorUInt32(d *DataUInt32, s *Synapses) byte {
	var (
		val uint32  //значение, которое анализируем
		f   float32 //применение функции анализа
		res byte    //количества вещества, которое хотим выплюнуть
	)
	//вычисляем величину медиатора
	if r.Typer%2 == 0 { //четный - положительный
		val = d.Data[r.NdataW]
	} else {
		val = 0xffffffff - d.Data[r.NdataW]
	}
	if r.Threshold > uint64(val) {
		return 0
	}
	f = float32(val) * float32(r.Force) / float32(r.Divforce)
	if f > 255 {
		f = 255
	} else if f < 1 {
		return 0
	}
	res = byte(f)

	//проходимся по всем аксонам и плюем эту величину
	for i := 0; i < 32; i++ {
		if r.A&(1<<i) != 0 { //проверяем, что данный аксон включен
			switch r.Typemedi {
			case NEURONACETILHOLIN: //ацетилхолин
				//вычисляем ячейку, куда плевать и вызываем метод аксона для соответсвующего выещества
				return r.Axons[i].DoAChMediate(&s.syn[r.Axons[i].N], res)
			}

		}
	}
	return 0
}
