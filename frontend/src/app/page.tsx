'use client';

import { Canvas } from "@react-three/fiber";
import { OrbitControls } from '@react-three/drei';
import Box from "@/components/box";
import { useState } from "react";


type User = {
  user_id: string;
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


export default function Home() {

  const [users, setUsers] = useState<User[]>(createUsers(100));

  return (
    <div className="font-sans">
      <main className="">
        <div className="bg-white w-screen h-screen">
          <Canvas>
            <OrbitControls />
            <ambientLight intensity={Math.PI / 2} />
            <spotLight position={[10, 10, 10]} angle={0.15} penumbra={1} decay={0} intensity={Math.PI} />
            <pointLight position={[-10, -10, -10]} decay={0} intensity={Math.PI} />
            {users.map((user: User) => (
              <Box key={user.user_id} active={user.active} position={[user.x, user.y, 0]}/>
            ))}
          </Canvas>
        </div>
      </main>
    </div>
  );
}
