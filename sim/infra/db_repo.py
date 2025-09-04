from dataclasses import asdict
from core.contracts import Impression, Action, Match
from typing import Optional
from core.repo import Repository
import httpx
import asyncio
from datetime import datetime


class DbRepo(Repository):
    def __init__(self, addr: str):
        self.addr = addr  # http://locahost:8080

    async def heartbeat(self, user_id: str, at: datetime):
        async with httpx.AsyncClient() as client:
            try:
                body = {
                    "active": True,
                    "at": at
                }
                r = await client.patch(f"{self.addr}/user/{user_id}/active",
                                       json=body)
                r.raise_for_status()
                print(r.json())
            except Exception as e:
                print(f"Exception | Heartbeat | {e}")

    async def add_impression(self, impression: Impression):
        async with httpx.AsyncClient() as client:
            try:
                r = await client.post(f"{self.addr}/impressions",
                                      data=impression)
                r.raise_for_status()
                print(r.json())
            except Exception as e:
                print(f"Exception | Add Impression | {e}")

    async def add_action(self, action: Action) -> Optional[Match]:
        mat = None
        async with httpx.AsyncClient() as client:
            try:
                r = await client.post(f"{self.addr}/actions",
                                      data=asdict(action))
                r.raise_for_status()
                mat = r.json()
            except Exception as e:
                print(f"Exception | Add Action | {e}")
        return mat


async def main():
    async with httpx.AsyncClient() as client:
        try:
            r = await client.get("http://localhost:8080/users/user1")
            print(r.json())
        except Exception as e:
            print(f"Error: {e}")

asyncio.run(main())
