package agent

//Effector - минимальная единица, передающая вещества в выходные данные
/*
Является, по сути переводчиком с языка веществ на язык цифр (информации)

Эффекторы могут менять значения в выходных файлах следующим образом:
у эффекторов вместо аксонов выходная величина гестеризует включение и выключение конкретного бита в выходном файле.
Ну или байта?
*/
type Effector struct {
}