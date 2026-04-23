import { Redirect } from 'expo-router';
import { useEffect, useState } from 'react';
import { getAccessToken } from '@/src/lib/auth';

export default function Index() {
  const [dest, setDest] = useState<'login' | 'tabs' | null>(null);

  useEffect(() => {
    getAccessToken().then(token => setDest(token ? 'tabs' : 'login'));
  }, []);

  if (!dest) return null;
  return <Redirect href={dest === 'tabs' ? '/(tabs)' : '/(auth)/login'} />;
}
