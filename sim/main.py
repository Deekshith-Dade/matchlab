import asyncio
import random
# from infra.memory_repo import InMemoryRepo
from infra.db_repo import DbRepo
from core.agent import UserAgent


def seed_positions(n=500, w=250, h=250, seed=7):
    r = random.Random(seed)
    return {f"u{i:04d}":
            (r.randint(int(-w/2), int(w/2)),  # x
            r.randint(int(-h/2), int(h/2)),  # y
            r.randint(3, 12))  # distance
            for i in range(n)}


async def main():
    users_data = seed_positions(n=400)
    # repo = InMemoryRepo(users_xy)
    url = "http://localhost:8080"
    repo = DbRepo(url)

    tasks = []
    for uid in users_data.keys():
        data = users_data[uid]
        agent = UserAgent(uid, data, repo)
        await repo.create_user(uid, data[0], data[1], data[2])
        t = asyncio.create_task(agent.run())
        tasks.append(t)
        await asyncio.sleep(5.0)

#    async def monitor():
#        while True:
#            m = await repo.match_q.get()
#            print(f"[monitor] {m.u1} MATCHED {m.u2} @ {m.at.isoformat()}")
#    mon = asyncio.create_task(monitor())

    await asyncio.gather(*tasks)

if __name__ == "__main__":
    asyncio.run(main())
