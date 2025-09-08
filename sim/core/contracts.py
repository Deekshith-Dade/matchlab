from dataclasses import dataclass
from typing import Literal
from datetime import datetime


Kind = Literal["like", "pass", "superlike"]


@dataclass
class Impression:
    viewer_id: str
    viewed_id: str
    rank: int
    at: datetime


@dataclass
class Action:
    viewer_id: str
    viewed_id: str
    kind: Kind
    at: datetime


@dataclass
class Match:
    u1: str
    u2: str
    at: datetime
