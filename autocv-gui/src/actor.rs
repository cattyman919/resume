use tokio::sync::mpsc;

pub enum ActorMessage {
    ReadFile(String),
    ModifyBuffer(String),
}

pub struct Actor {
    receiver: mpsc::Receiver<ActorMessage>,
}
