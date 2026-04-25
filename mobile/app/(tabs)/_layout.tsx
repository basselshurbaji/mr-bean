import { Tabs } from 'expo-router';
import { Feather } from '@expo/vector-icons';
import { colors } from '@/src/theme';
import { UserProvider } from '@/src/context/UserContext';

export default function TabsLayout() {
  return (
    <UserProvider>
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
        <Tabs.Screen
          name="index"
          options={{
            title: 'Home',
            tabBarIcon: ({ color, size }) => <Feather name="home" size={size} color={color} />,
          }}
        />
        <Tabs.Screen
          name="beans"
          options={{
            title: 'Beans',
            tabBarIcon: ({ color, size }) => <Feather name="coffee" size={size} color={color} />,
          }}
        />
        <Tabs.Screen
          name="gear"
          options={{
            title: 'My Gear',
            tabBarIcon: ({ color, size }) => <Feather name="tool" size={size} color={color} />,
          }}
        />
        <Tabs.Screen
          name="profile"
          options={{
            title: 'Profile',
            tabBarIcon: ({ color, size }) => <Feather name="user" size={size} color={color} />,
          }}
        />
      </Tabs>
    </UserProvider>
  );
}
