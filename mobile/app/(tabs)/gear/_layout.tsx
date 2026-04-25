import { Stack } from 'expo-router';
import { GearProvider } from '@/src/context/GearContext';

export default function GearLayout() {
  return (
    <GearProvider>
      <Stack screenOptions={{ headerShown: false, animation: 'slide_from_right' }}>
        <Stack.Screen name="index" />
        <Stack.Screen name="[id]" />
      </Stack>
    </GearProvider>
  );
}
