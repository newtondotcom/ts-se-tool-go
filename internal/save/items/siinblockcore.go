package items

// SiiNBlockCore is the base type for many save-game blocks in the C# code.
// It mirrors the C# SiiNBlockCore class from CustomClasses/Save/Items/SiiNBlockCore.cs.
//
// In the C# code, SiiNBlockCore provides a removeWritenBlock method that removes
// a nameless block from a global control list (FormMain.SiiNunitData.NamelessControlList).
// In Go, this is a no-op since we don't have the same UI structure, but the method
// is provided for API compatibility.
type SiiNBlockCore struct{}

// RemoveWritenBlock removes a written block from the nameless control list.
// In the C# code, this removes the block from FormMain.SiiNunitData.NamelessControlList.
// In Go, this is a no-op since we don't maintain a global control list.
// This method is provided for API compatibility with the C# codebase.
func (s *SiiNBlockCore) RemoveWritenBlock(name string) {
	// No-op in Go: we don't maintain a global nameless control list
	// like the C# FormMain.SiiNunitData.NamelessControlList
}

