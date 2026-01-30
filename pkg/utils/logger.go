package utils

var UI_Log_Chan = make(chan string, 100)

func TacticalLog(msg string) {
    select {
    case UI_Log_Chan <- msg:
    default:
        // Drop log if channel is full to prevent deadlocks
    }
}