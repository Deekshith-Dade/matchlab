import httpx
import asyncio


async def main():
    async with httpx.AsyncClient() as client:
        try:
            r = await client.get("http://localhost:8080/users/user1")
            print(r.json())
        except Exception as e:
            print(f"Error: {e}")

asyncio.run(main())
