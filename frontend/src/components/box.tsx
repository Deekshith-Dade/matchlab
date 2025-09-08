'use client';

import { useRef, useState } from "react";
import { ThreeElements, useFrame } from "@react-three/fiber";
import { Mesh } from "three";
import * as THREE from 'three';

type boxProps = ThreeElements['mesh'] & {
  active: boolean
}

type ringGeometry = ThreeElements['mesh'] & {
  rad: number
}

const Sphere = ({rad, ...props}: ringGeometry) => {

  const sphereRef = useRef<Mesh>(null!);
  const fullRadius = useRef<number>(rad)
  const [radius, setRadius] = useState<number>(0)


  useFrame((state, _) => {
    const t = state.clock.getElapsedTime();
    const base = fullRadius.current / 2;
    const newRadius = base + Math.sin(t * 2) * base; 
    setRadius(newRadius);
  })
  return (
    <group>
    <mesh 
      {...props}
      ref={sphereRef}>
      <circleGeometry args={[radius]} />
      <meshBasicMaterial 
        color={"blue"}
        transparent
        opacity={0.1}
        blending={THREE.AdditiveBlending}
        depthWrite={false}
        side={THREE.DoubleSide}
      />
    </mesh>
     <mesh 
      {...props}
      ref={sphereRef}>
      <ringGeometry args={[radius, radius + 0.1, 48]} />
      <meshStandardMaterial 
        color={"blue"}
        transparent
        opacity={0.25}
        blending={THREE.AdditiveBlending}
        depthWrite={false}
        side={THREE.DoubleSide}
      />
    </mesh>
 
    </group>
  );
};

export default function Box({active, ...props}: boxProps) {
  const meshRef = useRef<Mesh>(null!);

  useFrame((_, delta) => {
    if (meshRef.current) meshRef.current.rotation.x += delta;
  });

  if(!active) return;

  return (
    <group>
    <mesh
      {...props}
      ref={meshRef}
    >
      <boxGeometry args={[1, 1, 1]} />
      <meshStandardMaterial color={active ? "hotpink" : "orange"} />
    </mesh>
    <Sphere position={props.position} rad={10.0} />
    </group>
  );
}

