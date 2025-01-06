import asyncio
import time
from websockets.asyncio.client import connect
from datetime import datetime
import threading


async def player_a():
    async with connect("ws://localhost:8080/ws?playerID=playerA") as websocket:
        await websocket.send("Hello world!")

        ix = 0
        while ix < 10:
            await websocket.send(f"The clock is {datetime.now()}")

            message = await websocket.recv()
            print(f"Player A recieved: {message}")

            time.sleep(5)

def player_a_wrapper():
    asyncio.run(player_a())

async def player_b():
    async with connect("ws://localhost:8080/ws?playerID=playerB") as websocket:
        await websocket.send("Hello world!")

        ix = 0
        while ix < 10:
            await websocket.send(f"Jag bryr mig inte vad klockan är!")

            message = await websocket.recv()
            print(f"Player B recieved: {message}")

            time.sleep(5)

def player_b_wrapper():
    asyncio.run(player_b())


if __name__ == "__main__":
    player_a_thread = threading.Thread(target=player_a_wrapper)
    player_b_thread = threading.Thread(target=player_b_wrapper)

    player_a_thread.start()
    player_b_thread.start()

    player_a_thread.join()
    player_b_thread.join()

