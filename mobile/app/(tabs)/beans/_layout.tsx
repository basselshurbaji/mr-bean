import { Stack } from 'expo-router';
import { BeansProvider } from '@/src/context/BeansContext';

export default function BeansLayout() {
  return (
    <BeansProvider>
      <Stack screenOptions={{ headerShown: false, animation: 'slide_from_right' }}>
        <Stack.Screen name="index" />
        <Stack.Screen name="[id]" />
      </Stack>
    </BeansProvider>
  );
}
