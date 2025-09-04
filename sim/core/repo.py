from abc import ABC, abstractmethod
from datetime import datetime
from core.contracts import Impression, Action, Match
from typing import Sequence, Optional


class Repository(ABC):

    @abstractmethod
    async def create_user(self,
                          user_id: str,
                          x: int, y: int, dist: int) -> None: ...

    @abstractmethod
    async def heartbeat(self,
                        user_id: str,
                        at: datetime, active: bool) -> None: ...

    @abstractmethod
    async def add_impression(self, impression: Impression) -> None: ...

    @abstractmethod
    async def add_action(self, action: Action) -> Optional[Match]: ...

    @abstractmethod
    async def recent_matches(self, since: datetime) -> Sequence[Match]: ...

    @abstractmethod
    async def recommended_for(self,
                              user_id:
                              str, top_k: int) -> list[tuple[str, int]]: ...

    @abstractmethod
    async def mark_matched(self, u1: str, u2: str, at: datetime) -> Match: ...
