package expvarmon

// UI represents UI renderer.
type UI interface {
    Init(UIData) error
    Close()
    Update(UIData)
}
