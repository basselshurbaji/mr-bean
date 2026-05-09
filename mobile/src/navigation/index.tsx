import { useEffect } from 'react';
import { NavigationContainer } from '@react-navigation/native';
import { createNativeStackNavigator } from '@react-navigation/native-stack';
import { createBottomTabNavigator } from '@react-navigation/bottom-tabs';
import { Feather } from '@expo/vector-icons';
import { StatusBar } from 'expo-status-bar';
import * as SplashScreen from 'expo-splash-screen';
import { SafeAreaProvider } from 'react-native-safe-area-context';
import { KeyboardProvider } from 'react-native-keyboard-controller';
import {
  useFonts,
  PlayfairDisplay_400Regular,
  PlayfairDisplay_700Bold,
  PlayfairDisplay_700Bold_Italic,
  PlayfairDisplay_900Black,
} from '@expo-google-fonts/playfair-display';
import {
  DMSans_400Regular,
  DMSans_500Medium,
  DMSans_700Bold,
} from '@expo-google-fonts/dm-sans';
import {
  JetBrainsMono_400Regular,
  JetBrainsMono_500Medium,
  JetBrainsMono_600SemiBold,
} from '@expo-google-fonts/jetbrains-mono';
import { AuthProvider, useAuth } from '@/src/context/AuthContext';
import { UserProvider } from '@/src/context/UserContext';
import { GearProvider } from '@/src/context/GearContext';
import { BeansProvider } from '@/src/context/BeansContext';
import { ExtractionsProvider } from '@/src/context/ExtractionsContext';
import { colors } from '@/src/theme';

import LoginScreen from '@/app/(auth)/login';
import HomeScreen from '@/app/(tabs)/home/index';
import GearScreen from '@/app/(tabs)/gear/index';
import GearDetailScreen from '@/app/(tabs)/gear/[id]';
import BeanListScreen from '@/app/(tabs)/beans/index';
import BeanDetailScreen from '@/app/(tabs)/beans/[id]';
import ProfileScreen from '@/app/(tabs)/profile';

// ─── Param lists ──────────────────────────────────────────────────────────────

export type GearStackParamList = {
  GearList: undefined;
  GearDetail: { id: string };
};

export type BeansStackParamList = {
  BeanList: undefined;
  BeanDetail: { id: string };
};

// ─── Navigators ───────────────────────────────────────────────────────────────

const Root = createNativeStackNavigator();
const Auth = createNativeStackNavigator();
const Tab = createBottomTabNavigator();
const GearStack = createNativeStackNavigator<GearStackParamList>();
const BeansStack = createNativeStackNavigator<BeansStackParamList>();

function AuthNavigator() {
  return (
    <Auth.Navigator screenOptions={{ headerShown: false }}>
      <Auth.Screen name="Login" component={LoginScreen} />
    </Auth.Navigator>
  );
}

function GearNavigator() {
  return (
    <GearProvider>
      <GearStack.Navigator screenOptions={{ headerShown: false, animation: 'slide_from_right' }}>
        <GearStack.Screen name="GearList" component={GearScreen} />
        <GearStack.Screen name="GearDetail" component={GearDetailScreen} />
      </GearStack.Navigator>
    </GearProvider>
  );
}

function BeansNavigator() {
  return (
    <BeansProvider>
      <BeansStack.Navigator screenOptions={{ headerShown: false, animation: 'slide_from_right' }}>
        <BeansStack.Screen name="BeanList" component={BeanListScreen} />
        <BeansStack.Screen name="BeanDetail" component={BeanDetailScreen} />
      </BeansStack.Navigator>
    </BeansProvider>
  );
}

function HomeTab() {
  return (
    <BeansProvider>
      <GearProvider>
        <ExtractionsProvider>
          <HomeScreen />
        </ExtractionsProvider>
      </GearProvider>
    </BeansProvider>
  );
}

function MainTabs() {
  return (
    <UserProvider>
      <Tab.Navigator
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
        <Tab.Screen
          name="Home"
          component={HomeTab}
          options={{
            title: 'Home',
            tabBarIcon: ({ color, size }) => <Feather name="home" size={size} color={color} />,
          }}
        />
        <Tab.Screen
          name="Beans"
          component={BeansNavigator}
          options={{
            title: 'Beans',
            tabBarIcon: ({ color, size }) => <Feather name="coffee" size={size} color={color} />,
          }}
        />
        <Tab.Screen
          name="Gear"
          component={GearNavigator}
          options={{
            title: 'My Gear',
            tabBarIcon: ({ color, size }) => <Feather name="tool" size={size} color={color} />,
          }}
        />
        <Tab.Screen
          name="Profile"
          component={ProfileScreen}
          options={{
            title: 'Profile',
            tabBarIcon: ({ color, size }) => <Feather name="user" size={size} color={color} />,
          }}
        />
      </Tab.Navigator>
    </UserProvider>
  );
}

function AppInner() {
  const { isAuthenticated, ready } = useAuth();
  const [fontsLoaded, fontError] = useFonts({
    PlayfairDisplay_400Regular,
    PlayfairDisplay_700Bold,
    PlayfairDisplay_700Bold_Italic,
    PlayfairDisplay_900Black,
    DMSans_400Regular,
    DMSans_500Medium,
    DMSans_700Bold,
    JetBrainsMono_400Regular,
    JetBrainsMono_500Medium,
    JetBrainsMono_600SemiBold,
  });

  useEffect(() => {
    if (fontError) throw fontError;
  }, [fontError]);

  useEffect(() => {
    if (fontsLoaded && ready) SplashScreen.hideAsync();
  }, [fontsLoaded, ready]);

  if (!fontsLoaded || !ready) return null;

  return (
    <NavigationContainer>
      <Root.Navigator screenOptions={{ headerShown: false }}>
        {isAuthenticated ? (
          <Root.Screen name="Main" component={MainTabs} />
        ) : (
          <Root.Screen name="Auth" component={AuthNavigator} />
        )}
      </Root.Navigator>
    </NavigationContainer>
  );
}

SplashScreen.preventAutoHideAsync();

export default function App() {
  return (
    <SafeAreaProvider>
      <KeyboardProvider>
        <StatusBar style="dark" />
        <AuthProvider>
          <AppInner />
        </AuthProvider>
      </KeyboardProvider>
    </SafeAreaProvider>
  );
}