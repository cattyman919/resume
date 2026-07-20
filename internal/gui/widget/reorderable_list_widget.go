package widget

import (
	"fmt"
	"unsafe"

	"github.com/AllenDang/cimgui-go/imgui"
	"github.com/AllenDang/giu"
)

type ReorderableListWidget[T any] struct {
	id            string
	layout        []T
	layoutDisplay giu.Layout
	onReorder     func([]T)
}

func ReorderableList[T any](id string, layout []T) *ReorderableListWidget[T] {
	return &ReorderableListWidget[T]{
		id:     id,
		layout: layout,
	}
}

func (r *ReorderableListWidget[T]) Layout(layoutDisplay ...giu.Widget) *ReorderableListWidget[T] {
	r.layoutDisplay = layoutDisplay
	return r
}

func (r *ReorderableListWidget[T]) OnReorder(onReorder func([]T)) *ReorderableListWidget[T] {
	r.onReorder = onReorder
	return r
}

func (r *ReorderableListWidget[T]) Build() {
	if r.id == "" {
		panic("ReorderableList: id is required")
	}

	if r.onReorder == nil {
		panic("ReorderableList: OnReorder callback is required")
	}

	var layoutDisplay giu.Layout
	r.layoutDisplay.Range(func(w giu.Widget) {
		layoutDisplay = append(layoutDisplay, w)
	})

	if len(r.layout) != len(layoutDisplay) {
		return
	}

	clipper := imgui.NewListClipper()
	defer clipper.Destroy()

	clipper.Begin(int32(len(layoutDisplay)))

	var sourceIdx, targetIdx int32 = -1, -1

	for clipper.Step() {
		for i := clipper.DisplayStart(); i < clipper.DisplayEnd(); i++ {
			r.buildItem(int(i), layoutDisplay[i], &sourceIdx, &targetIdx)
		}
	}

	clipper.End()

	if sourceIdx >= 0 && targetIdx >= 0 && sourceIdx != targetIdx {
		newLayout := r.reordered(int(sourceIdx), int(targetIdx))
		r.onReorder(newLayout)
	}
}

func (r *ReorderableListWidget[T]) buildItem(i int, display giu.Widget, sourceIdx, targetIdx *int32) {
	startPos := imgui.CursorPos()
	width := imgui.ContentRegionAvail().X

	// Draw the caller-provided content first.
	display.Build()

	// Draw a full-width selectable over the content to make the whole row
	// interactive and provide an item ID for drag-and-drop.
	rowHeight := imgui.ItemRectSize().Y
	imgui.SetCursorPos(startPos)
	imgui.SelectableBoolV(
		fmt.Sprintf("##%s-row%d", r.id, i),
		false,
		imgui.SelectableFlagsNone,
		imgui.Vec2{X: width, Y: rowHeight},
	)

	rowMin := imgui.ItemRectMin()
	rowMax := imgui.ItemRectMax()

	// Drag source.
	if imgui.BeginDragDropSource() {
		idx := int32(i)
		imgui.SetDragDropPayload(
			r.payloadType(),
			uintptr(unsafe.Pointer(&idx)),
			uint64(unsafe.Sizeof(idx)),
		)
		display.Build()
		imgui.EndDragDropSource()
	}

	// Drop target.
	if imgui.BeginDragDropTarget() {
		payload := imgui.AcceptDragDropPayload(r.payloadType())
		if payload != nil && payload.CData != nil && payload.DataSize() == int32(unsafe.Sizeof(int32(0))) {
			src := *(*int32)(unsafe.Pointer(payload.Data()))
			if src != int32(i) {
				*sourceIdx = src
				*targetIdx = int32(i)

				// Highlight the row.
				col := imgui.ColorU32ColV(imgui.ColDragDropTargetBg, 0.5)
				imgui.WindowDrawList().AddRectFilled(rowMin, rowMax, col)
			}
		}

		imgui.EndDragDropTarget()
	}
}

func (r *ReorderableListWidget[T]) reordered(source, target int) []T {
	if source == target {
		newLayout := make([]T, len(r.layout))
		copy(newLayout, r.layout)
		return newLayout
	}

	newLayout := make([]T, len(r.layout))
	copy(newLayout, r.layout)

	item := newLayout[source]
	if source < target {
		copy(newLayout[source:], newLayout[source+1:target])
		newLayout[target-1] = item
	} else {
		copy(newLayout[target+1:source+1], newLayout[target:source])
		newLayout[target] = item
	}

	return newLayout
}

func (r *ReorderableListWidget[T]) payloadType() string {
	if len(r.id) > 31 {
		return r.id[:31]
	}
	return r.id
}
