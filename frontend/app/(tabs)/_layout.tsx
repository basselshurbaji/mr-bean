import { Tabs } from 'expo-router';
import { colors } from '@/src/theme';

export default function TabsLayout() {
  return (
    <Tabs
      screenOptions={{
        headerShown: false,
        tabBarActiveTintColor: colors.fgPrimary,
        tabBarInactiveTintColor: colors.fgTertiary,
        tabBarStyle: {
          backgroundColor: colors.bgApp,
          borderTopColor: colors.borderSubtle,
        },
      }}
    >
      <Tabs.Screen name="index" options={{ title: 'Home' }} />
      <Tabs.Screen name="beans" options={{ title: 'Beans' }} />
    </Tabs>
  );
}
