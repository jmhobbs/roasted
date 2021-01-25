package sr700

type Program []State

type State struct {
	Fan   Speed
	Timer Timer
	Heat  Heat
}
