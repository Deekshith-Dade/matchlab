from datetime import datetime, timezone
import asyncio
import random
from core.contracts import Impression, Action
from core.repo import Repository


def now():
    return datetime.now(timezone.utc)


def like_probability(rank: int) -> float:
    return max(0.05, 0.35 - 0.01 * rank)


class UserAgent:
    def __init__(self, user_id: str, data: tuple[int, int, int],
                 repo: Repository, *, jitter_ms=(50, 250)):
        self.user_id = user_id
        self.user_x = data[0]
        self.user_y = data[1]
        self.distance = data[2]
        self.repo = repo
        self.jitter_ms = jitter_ms
        self.rand = random.Random(hash(user_id) & 0xffffffff)
        print(f"User: {user_id} logged in")

    async def run(self):
        while True:
            await self.repo.heartbeat(self.user_id, now(), active=True)

            recs = await self.repo.recommended_for(self.user_id, 5)
            for viewed_id, rank in recs:
                t = now()
                await self.repo.add_impression(Impression(
                    viewer_id=self.user_id, viewed_id=viewed_id,
                    rank=rank, at=t
                ))
                if self.rand.random() < like_probability(rank):
                    match = await self.repo.add_action(Action(
                        viewer_id=self.user_id, viewed_id=viewed_id,
                        kind="LIKE", at=now()
                    ))
                    print(f"User: {self.user_id} LIKED {viewed_id}")
                    if match:
                        # User with some probability doesn't like his match and goes back to dating
                        if self.rand.random() < 0.75:
                            print(f"User {self.user_id} Doesn't like his match")
                            continue
                        print(f"Users {match.u1} LOVEEES {match.u2}")
                        return
            ms = self.rand.uniform(*self.jitter_ms)/10.0
            await self.repo.heartbeat(self.user_id, now(), active=False)
            await asyncio.sleep(ms)
