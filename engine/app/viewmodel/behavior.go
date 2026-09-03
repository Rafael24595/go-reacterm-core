package viewmodel

// BehaviorContext holds operational flags that instruct the engine
// on how to handle the execution cycle for the active screen.
type BehaviorContext struct {
    // NeedsPulse indicates whether the engine should continue to invoke the screen's pulse method.
    NeedsPulse bool 
}
