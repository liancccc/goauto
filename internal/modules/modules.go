package modules

var Modules []Module

func RegisterModule(module Module) {
	Modules = append(Modules, module)
}
