package models

//Struct to define ClickUp event
type ClickUpEvent struct {
    Event string `json:"event"`
    Task  struct {
        Priority string `json:"priority"`
    } `json:"task"`
}
