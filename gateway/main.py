from fastapi import FastAPI
from controller.v1.routes.chat import router as chat_router
import uvicorn
from repository.sqlite import create_tables



app = FastAPI()

app.include_router(chat_router, prefix="/v1")

if __name__ == "__main__":
    create_tables()
    uvicorn.run(app, host="0.0.0.0", port=8080)
