'use client';

import { Canvas } from "@react-three/fiber";
import { OrbitControls } from '@react-three/drei';
import Box from "@/components/box";
import { useEffect, useState } from "react";


type User = {
  id: string;
  x: number;
  y: number;
  active: boolean;
  distance: number; //radial distance to which he's okay looking for
}

const createUsers = (n: number): User[] => {
  const spread = 200;
  const users: User[] = [];
  for (let i = 0; i < n; i++) {
    const x = Math.random() * spread - spread/2;
    const y = Math.random() * spread - spread/2;
    const active = Math.random() > 0.5 ? false : true;
    users.push({ user_id: `user${i}`, x: x, y: y, active: active, distance: 50 })
  }
  return users;
}

const getUsers = async () : Promise<User[]> => {
  try{
  const res = await fetch("http://localhost:8080/users");
    if(!res.ok){
      const err = await res.json()
      throw new Error(`HTTP error! Status: ${res.status}, Message: ${err}`)
    }
  const data = await res.json()
  return data;
  }
  catch(error){
    console.log(error);
    return [];
  }
  }


export default function Home() {

  const [users, setUsers] = useState<User[]>([]);
  
  useEffect(() => {
    let mounted = true;

    getUsers().then((data) => {
      if(mounted) setUsers(data)
    });

    const interval = setInterval(() => {
        getUsers().then((data: User[]) => {
          if(mounted) setUsers(data);
        }
        )
    }, 500);

    return () => {
      mounted = false;
      clearInterval(interval);
    }
  }, []);

  // console.log(users)

  return (
    <div className="font-sans">
      <main className="">
        <div className="bg-white w-screen h-screen">
          <Canvas>
            <OrbitControls />
            <ambientLight intensity={Math.PI / 2} />
            <spotLight position={[10, 10, 10]} angle={0.15} penumbra={1} decay={0} intensity={Math.PI} />
            <pointLight position={[-10, -10, -10]} decay={0} intensity={Math.PI} />
            {users && users.map((user: User) => (
              <Box key={user.id} active={user.active} position={[user.x, user.y, -30]}/>
            ))}
          </Canvas>
        </div>
      </main>
    </div>
  );
}
