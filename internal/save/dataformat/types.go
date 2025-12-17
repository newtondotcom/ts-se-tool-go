package dataformat

// This package contains low-level data types used across save-game items,
// mirroring CustomClasses/Save/DataFormat in the C# project.

// Color approximates SCS_Color from the C# code. It stores channels as bytes.
type Color struct {
	A uint8
	R uint8
	G uint8
	B uint8
}

// Float is a thin wrapper around float32 mirroring SCS_Float. It allows us to
// add helper methods later if needed.
type Float float32

// String wraps a plain string; added for symmetry with SCS_string.
type String string

// Vector2f mirrors Vector_2f.
type Vector2f struct {
	X float32
	Y float32
}

// Vector3f mirrors Vector_3f.
type Vector3f struct {
	X float32
	Y float32
	Z float32
}

// Vector4f mirrors Vector_4f and the 3f_4f variants.
type Vector4f struct {
	X float32
	Y float32
	Z float32
	W float32
}

// Vector3i mirrors Vector_3i.
type Vector3i struct {
	X int32
	Y int32
	Z int32
}


