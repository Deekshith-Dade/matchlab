import asyncio
from datetime import datetime
from typing import Optional, Sequence
from core.contracts import Impression, Action, Match
from core.repo import Repository


class InMemoryRepo(Repository):
    def __init__(self, users_xy: dict[str, tuple[int, int]]):
        self._lock = asyncio.Lock()
        self.users_xy = users_xy
        self.last_seen: dict[str, datetime] = {}
        self.active: set[str] = set()
        self.likes: set[tuple[str, str]] = set()
        self.matched:  set[frozenset[str]] = set()
        self.impressions: list[Impression] = []
        self.actions: list[Action] = []
        self.matches: list[Match] = []
        self.match_q: asyncio.Queue[Match] = asyncio.Queue()

    async def heartbeat(self, user_id: str, at: datetime):
        async with self._lock:
            self.last_seen[user_id] = at
            self.active.add(user_id)

    async def add_impression(self, impression: Impression):
        async with self._lock:
            self.impressions.append(impression)

    async def add_action(self, action: Action) -> Optional[Match]:
        mat = None
        async with self._lock:
            self.actions.append(action)

            if action.kind == "LIKE":
                self.likes.add((action.viewer_id, action.viewed_id))

                if (action.viewed_id, action.viewer_id) in self.likes:
                    key = frozenset({action.viewer_id, action.viewed_id})
                    if key not in self.matched:
                        mat = Match(u1=action.viewer_id,
                                    u2=action.viewed_id,
                                    at=action.at)
                        self.matched.add(key)
                        self.matches.append(mat)
                        self.active.discard(action.viewer_id)
                        self.active.discard(action.viewed_id)
        if mat:
            await self.match_q.put(mat)
        return mat

    async def recent_matches(self, since: datetime) -> Sequence[Match]:
        async with self._lock:
            return [m for m in self.matches if m.at >= since]

    async def mark_matched(self, u1: str, u2: str, at: datetime) -> Match:
        async with self._lock:
            key = frozenset({u1, u2})
            if key not in self.matched:
                mat = Match(u1, u2, at)
                self.matched.add(key)
                self.matches.append(mat)
                self.active.discard(u1)
                self.active.discard(u2)
                return mat
            else:
                return next(m for m in self.matches
                            if frozenset({m.u1, m.u2}) == key)

    async def recommended_for(self, user_id: str,
                              top_k: int) -> list[tuple[str, int]]:
        async with self._lock:
            ux, uy = self.users_xy[user_id]

            def dist(v):
                x, y = self.users_xy[v]
                return abs(x-ux) + abs(y-uy)
            candidates = [
                v for v in self.users_xy
                if v != user_id and
                v in self.active and
                (user_id, v) not in self.likes
            ]
            candidates.sort(key=dist)
            return [(v, r+1) for r, v in enumerate(candidates[:top_k])]
