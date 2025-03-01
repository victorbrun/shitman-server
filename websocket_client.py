import asyncio
import time
import json
from websockets.asyncio.client import connect
from datetime import datetime
import threading


async def player_a():
    command = {
        "player_id": "playerA",
        "play_cards": [
            {"rank": "Seven", "suit": "Hearts"},
            {"rank": "King", "suit": "Spades"}
        ],
        "play_card_from_hidden_hand": True,
        "play_random_card_from_deck": False
    }

    async with connect("ws://localhost:8080/ws?player_id=playerA") as websocket:
        ix = 0
        while ix < 10:
            await websocket.send(json.dumps(command))

            message = await websocket.recv()
            print(f"Player A recieved: {message}")

            time.sleep(5)

def player_a_wrapper():
    asyncio.run(player_a())

async def player_b():
    command = {
        "player_id": "playerB",
        "play_cards": [
            {"rank": "Ace", "suit": "Hearts"},
            {"rank": "King", "suit": "Spades"}
        ],
        "play_card_from_hidden_hand": False,
        "play_random_card_from_deck": False
    }

    async with connect("ws://localhost:8080/ws?player_id=playerB&game_id=1") as websocket:
        ix = 0
        while ix < 10:
            await websocket.send(json.dumps(command))

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

