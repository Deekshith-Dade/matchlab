import asyncio
from core.contracts import Impression, Action, Match
from typing import Optional, Sequence
from core.repo import Repository
import httpx
from datetime import datetime, timezone


class DbRepo(Repository):
    def __init__(self, addr: str):
        self.addr = addr  # http://locahost:8080

    async def create_user(self,
                          user_id: str,
                          x: int, y: int, dist: int) -> None:
        async with httpx.AsyncClient() as client:
            try:
                data = {
                    "id": user_id,
                    "x": x,
                    "y": y,
                    "distance": dist,
                    "active": False
                }
                r = await client.post(f"{self.addr}/users", json=data)
                r.raise_for_status()
            except Exception as e:
                print(f"Exception | Create User | {e}")

    async def heartbeat(self, user_id: str, at: datetime, active: bool):
        async with httpx.AsyncClient() as client:
            try:
                body = {
                    "active": active,
                    "at": at.isoformat()
                }
                r = await client.patch(f"{self.addr}/users/{user_id}/active",
                                       json=body)
                r.raise_for_status()
            except Exception as e:
                print(f"Exception | Heartbeat | {e}")

    async def add_impression(self, impression: Impression):
        async with httpx.AsyncClient() as client:
            try:
                data = {
                    "viewer_id": impression.viewer_id,
                    "viewed_id": impression.viewed_id,
                    "rank": impression.rank
                }
                r = await client.post(f"{self.addr}/impressions",
                                      json=data)
                r.raise_for_status()
                print("Imperssion created", r.status_code)
            except Exception as e:
                print(f"Exception | Add Impression | {e}")

    async def add_action(self, action: Action) -> Optional[Match]:
        mat = None
        async with httpx.AsyncClient() as client:
            try:
                data = {
                    "viewer_id": action.viewer_id,
                    "viewed_id": action.viewed_id,
                    "kind": action.kind,
                    "at": action.at.isoformat()
                }
                r = await client.post(f"{self.addr}/actions",
                                      json=data)
                r.raise_for_status()
                mat = r.json()
            except Exception as e:
                print(f"Exception | Add Action | {e}")
        return mat

    async def recent_matches(self, since: datetime) -> Sequence[Match]:
        print(since)
        return []

    async def mark_matched(self, u1: str, u2: str, at: datetime) -> Match:
        return Match(
            u1=u1,
            u2=u2,
            at=at
        )

    async def recommended_for(self, user_id: str,
                              top_k: int) -> list[tuple[str, int]]:
        async with httpx.AsyncClient() as client:
            out = []
            try:
                r = await client.get(
                    f"{self.addr}/recommendations/{user_id}?topk={top_k}")
                r.raise_for_status()
                data = r.json()
                if data:
                    out = [(entry['user_id'], entry['rank']) for entry in data]
            except Exception as e:
                print(f"Exception | Recommendation | {e}")
        return out


async def main():
    url = "http://localhost:8080"
    repo = DbRepo(url)
    out = await repo.recommended_for("user3", 2)
    print(out)

    action = Action(
        viewer_id="user1",
        viewed_id="user4",
        kind="like",
        at=datetime.now(timezone.utc)
    )
    data = await repo.add_action(action)
    print(f"Add action respnse data {data}")


if __name__ == "__main__":
    asyncio.run(main())
