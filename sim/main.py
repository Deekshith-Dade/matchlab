import asyncio
import random
from infra.memory_repo import InMemoryRepo
from core.agent import UserAgent


def seed_positions(n=500, w=1000, h=1000, seed=7):
    r = random.Random(seed)
    return {f"u{i:04d}": (r.randint(20, w-20), r.randint(20, h-20))
            for i in range(n)}


async def main():
    users_xy = seed_positions(n=400)
    repo = InMemoryRepo(users_xy)

    tasks = []
    for uid in users_xy.keys():
        agent = UserAgent(uid, repo)
        t = asyncio.create_task(agent.run())
        tasks.append(t)
        await asyncio.sleep(1.0)

    async def monitor():
        while True:
            m = await repo.match_q.get()
            print(f"[monitor] {m.u1} MATCHED {m.u2} @ {m.at.isoformat()}")
    mon = asyncio.create_task(monitor())

    await asyncio.gather(mon, *tasks)

if __name__ == "__main__":
    asyncio.run(main())
