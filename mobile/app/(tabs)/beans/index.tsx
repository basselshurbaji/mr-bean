import {
  Animated,
  FlatList,
  Pressable,
  StyleSheet,
  Text,
  View,
} from 'react-native';
import { useCallback, useRef, useState } from 'react';
import { useNavigation } from '@react-navigation/native';
import { NativeStackNavigationProp } from '@react-navigation/native-stack';
import { BeansStackParamList } from '@/src/navigation';
import { useSafeAreaInsets } from 'react-native-safe-area-context';
import { Feather } from '@expo/vector-icons';
import { colors, palette, radii, spacing } from '@/src/theme';
import { useBeans } from '@/src/context/BeansContext';
import { Bean, processLabel, roastLabel } from '@/src/api/beans';
import RoastBubble from '@/src/components/RoastBubble';
import BeanSheet from './BeanSheet';

// ─── Toast ────────────────────────────────────────────────────────────────────

function useToast() {
  const opacity   = useRef(new Animated.Value(0)).current;
  const translateY = useRef(new Animated.Value(40)).current;
  const [msg, setMsg] = useState('');
  const timer = useRef<ReturnType<typeof setTimeout> | null>(null);

  function show(text: string) {
    setMsg(text);
    if (timer.current) clearTimeout(timer.current);
    Animated.parallel([
      Animated.timing(opacity,    { toValue: 1, duration: 220, useNativeDriver: true }),
      Animated.timing(translateY, { toValue: 0, duration: 220, useNativeDriver: true }),
    ]).start();
    timer.current = setTimeout(() => {
      Animated.parallel([
        Animated.timing(opacity,    { toValue: 0, duration: 200, useNativeDriver: true }),
        Animated.timing(translateY, { toValue: 40, duration: 200, useNativeDriver: true }),
      ]).start();
    }, 2600);
  }

  return { show, msg, animatedStyle: { opacity, transform: [{ translateY }] } };
}

// ─── Bean card ────────────────────────────────────────────────────────────────

function BeanCard({ bean, onPress }: { bean: Bean; onPress: () => void }) {
  const scale = useRef(new Animated.Value(1)).current;
  const sub   = [bean.roaster, bean.origin].filter(Boolean).join(' · ');
  const proc  = processLabel(bean.process);
  const roast = roastLabel(bean.roast_level);

  return (
    <Pressable
      onPressIn={() =>
        Animated.spring(scale, { toValue: 0.98, useNativeDriver: true, speed: 30 }).start()
      }
      onPressOut={() =>
        Animated.spring(scale, { toValue: 1, useNativeDriver: true, speed: 30 }).start()
      }
      onPress={onPress}
    >
      <Animated.View style={[styles.card, { transform: [{ scale }] }]}>
        <RoastBubble roastId={bean.roast_level} size={50} />

        <View style={styles.cardText}>
          <Text style={styles.cardName} numberOfLines={1}>{bean.name}</Text>
          {sub ? <Text style={styles.cardSub} numberOfLines={1}>{sub}</Text> : null}
          {bean.tasting_notes ? (
            <Text style={styles.cardTasting} numberOfLines={1}>{bean.tasting_notes}</Text>
          ) : null}
        </View>

        <View style={styles.cardMeta}>
          {proc ? (
            <View style={styles.processBadge}>
              <Text style={styles.processBadgeText}>{proc}</Text>
            </View>
          ) : null}
          {roast ? <Text style={styles.roastLabel}>{roast}</Text> : null}
        </View>
      </Animated.View>
    </Pressable>
  );
}

// ─── Screen ───────────────────────────────────────────────────────────────────

export default function BeansScreen() {
  const insets = useSafeAreaInsets();
  const navigation = useNavigation<NativeStackNavigationProp<BeansStackParamList>>();
  const { beans, addBean, refresh } = useBeans();
  const toast = useToast();

  const [refreshing, setRefreshing] = useState(false);
  const handleRefresh = useCallback(async () => {
    setRefreshing(true);
    await refresh();
    setRefreshing(false);
  }, [refresh]);

  const [sheetOpen, setSheetOpen] = useState(false);

  const fabScale = useRef(new Animated.Value(1)).current;

  return (
    <View style={styles.screen}>
      {/* Header */}
      <View style={[styles.header, { paddingTop: insets.top + 12 }]}>
        <Text style={styles.title}>Your beans</Text>
        <Text style={styles.subtitle}>
          {beans.length} bean{beans.length !== 1 ? 's' : ''} in your library
        </Text>
      </View>

      {/* List */}
      <FlatList
        data={beans}
        keyExtractor={b => b.id}
        contentContainerStyle={styles.list}
        showsVerticalScrollIndicator={false}
        refreshing={refreshing}
        onRefresh={handleRefresh}
        ListEmptyComponent={
          <View style={styles.empty}>
            <Text style={styles.emptyEmoji}>🫘</Text>
            <Text style={styles.emptyTitle}>No beans on deck.</Text>
            <Text style={styles.emptySub}>Time to fix that, don't you think?</Text>
          </View>
        }
        renderItem={({ item }) => (
          <BeanCard
            bean={item}
            // eslint-disable-next-line @typescript-eslint/no-explicit-any
            onPress={() => navigation.navigate('BeanDetail', { id: item.id })}
          />
        )}
      />

      {/* FAB */}
      <Pressable
        style={styles.fab}
        onPressIn={() =>
          Animated.spring(fabScale, { toValue: 0.93, useNativeDriver: true, damping: 8, stiffness: 180 }).start()
        }
        onPressOut={() =>
          Animated.spring(fabScale, { toValue: 1, useNativeDriver: true, damping: 8, stiffness: 180 }).start()
        }
        onPress={() => setSheetOpen(true)}
      >
        <Animated.View style={[styles.fabInner, { transform: [{ scale: fabScale }] }]}>
          <Feather name="plus" size={24} color={palette.cream100} />
        </Animated.View>
      </Pressable>

      {/* Sheet */}
      {sheetOpen && (
        <BeanSheet
          onClose={() => setSheetOpen(false)}
          onSaved={bean => {
            addBean(bean);
            setSheetOpen(false);
            toast.show(`${bean.name} added ✓`);
          }}
        />
      )}

      {/* Toast */}
      <Animated.View style={[styles.toast, toast.animatedStyle]} pointerEvents="none">
        <Text style={styles.toastText}>{toast.msg}</Text>
      </Animated.View>
    </View>
  );
}

// ─── Styles ───────────────────────────────────────────────────────────────────

const styles = StyleSheet.create({
  screen: { flex: 1, backgroundColor: colors.bgApp },

  header: {
    paddingHorizontal: spacing[6],
    paddingBottom: spacing[3],
  },
  title: {
    fontFamily: 'PlayfairDisplay_900Black',
    fontSize: 34,
    lineHeight: 40,
    letterSpacing: -0.7,
    color: colors.fgPrimary,
  },
  subtitle: {
    fontFamily: 'DMSans_400Regular',
    fontSize: 13,
    color: palette.espresso400,
    marginTop: 2,
  },

  list: {
    paddingHorizontal: spacing[5],
    paddingBottom: 140,
    gap: 10,
  },

  card: {
    backgroundColor: palette.cream200,
    borderRadius: radii.lg,
    paddingVertical: 14,
    paddingHorizontal: 16,
    flexDirection: 'row',
    alignItems: 'center',
    gap: 12,
    shadowColor: palette.espresso800,
    shadowOffset: { width: 0, height: 2 },
    shadowOpacity: 0.07,
    shadowRadius: 8,
    elevation: 2,
  },
  cardText: { flex: 1, minWidth: 0 },
  cardName: {
    fontFamily: 'DMSans_700Bold',
    fontSize: 15,
    color: colors.fgPrimary,
  },
  cardSub: {
    fontFamily: 'DMSans_400Regular',
    fontSize: 12,
    color: palette.espresso400,
    marginTop: 2,
  },
  cardTasting: {
    fontFamily: 'DMSans_400Regular',
    fontSize: 12,
    color: palette.espresso300,
    marginTop: 1,
  },
  cardMeta: { alignItems: 'flex-end', gap: 4, flexShrink: 0 },
  processBadge: {
    backgroundColor: palette.cream300,
    paddingHorizontal: 8,
    paddingVertical: 3,
    borderRadius: radii.full,
  },
  processBadgeText: {
    fontFamily: 'DMSans_700Bold',
    fontSize: 10,
    color: palette.espresso600,
  },
  roastLabel: {
    fontFamily: 'DMSans_400Regular',
    fontSize: 11,
    color: palette.cream600,
  },

  empty: {
    alignItems: 'center',
    paddingTop: 52,
    paddingHorizontal: spacing[6],
  },
  emptyEmoji:  { fontSize: 40, marginBottom: 12 },
  emptyTitle: {
    fontFamily: 'DMSans_700Bold',
    fontSize: 16,
    color: colors.fgPrimary,
    marginBottom: 6,
  },
  emptySub: {
    fontFamily: 'DMSans_400Regular',
    fontSize: 13,
    color: palette.espresso300,
    textAlign: 'center',
  },

  fab: {
    position: 'absolute',
    bottom: 12,
    right: 22,
  },
  fabInner: {
    width: 56,
    height: 56,
    borderRadius: 28,
    backgroundColor: palette.espresso800,
    alignItems: 'center',
    justifyContent: 'center',
    shadowColor: palette.espresso800,
    shadowOffset: { width: 0, height: 6 },
    shadowOpacity: 0.38,
    shadowRadius: 20,
    elevation: 8,
  },

  toast: {
    position: 'absolute',
    bottom: 104,
    left: spacing[5],
    right: spacing[5],
    backgroundColor: colors.fgPrimary,
    borderRadius: 16,
    paddingVertical: 14,
    paddingHorizontal: 18,
  },
  toastText: {
    fontFamily: 'DMSans_500Medium',
    fontSize: 14,
    color: colors.fgInverse,
    textAlign: 'center',
  },
});
