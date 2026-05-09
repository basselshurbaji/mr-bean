import {
  ActivityIndicator,
  Animated,
  Pressable,
  ScrollView,
  StyleSheet,
  Text,
  View,
} from 'react-native';
import { useRoute, useNavigation } from '@react-navigation/native';
import { RouteProp } from '@react-navigation/native';
import { NativeStackNavigationProp } from '@react-navigation/native-stack';
import { GearStackParamList } from '@/src/navigation';
import { useRef, useState } from 'react';
import { useSafeAreaInsets } from 'react-native-safe-area-context';
import { Feather } from '@expo/vector-icons';
import { colors, palette, radii, spacing } from '@/src/theme';
import { useGear } from '@/src/context/GearContext';
import { gearApi } from '@/src/api/gear';
import GearIcon from '@/src/components/GearIcon';
import GearSheet from './GearSheet';

const TYPES: Record<string, string> = {
  machine: 'Espresso machine',
  grinder: 'Grinder',
  scale: 'Scale',
  portafilter: 'Portafilter',
  tamper: 'Tamper',
  distributor: 'Distribution tool',
  wdt: 'WDT tool',
  basket: 'Basket',
  puckscreen: 'Puck screen',
  other: 'Other',
};

function useToast() {
  const opacity = useRef(new Animated.Value(0)).current;
  const translateY = useRef(new Animated.Value(40)).current;
  const [msg, setMsg] = useState('');
  const timer = useRef<ReturnType<typeof setTimeout> | null>(null);

  function show(text: string) {
    setMsg(text);
    if (timer.current) clearTimeout(timer.current);
    Animated.parallel([
      Animated.timing(opacity, { toValue: 1, duration: 220, useNativeDriver: true }),
      Animated.timing(translateY, { toValue: 0, duration: 220, useNativeDriver: true }),
    ]).start();
    timer.current = setTimeout(() => {
      Animated.parallel([
        Animated.timing(opacity, { toValue: 0, duration: 200, useNativeDriver: true }),
        Animated.timing(translateY, { toValue: 40, duration: 200, useNativeDriver: true }),
      ]).start();
    }, 2600);
  }

  return { show, msg, animatedStyle: { opacity, transform: [{ translateY }] } };
}

export default function GearDetailScreen() {
  const route = useRoute<RouteProp<GearStackParamList, 'GearDetail'>>();
  const { id } = route.params;
  const navigation = useNavigation<NativeStackNavigationProp<GearStackParamList>>();
  const insets = useSafeAreaInsets();
  const { gear, updateGear, removeGear } = useGear();
  const toast = useToast();

  const item = gear.find(g => g.id === id);

  const [editOpen, setEditOpen] = useState(false);
  const [deleting, setDeleting] = useState(false);

  async function handleDelete() {
    if (deleting || !item) return;
    setDeleting(true);
    try {
      await gearApi.deleteGear(item.id);
      removeGear(item.id);
      navigation.goBack();
    } catch {
      setDeleting(false);
      toast.show('Failed to remove item.');
    }
  }

  if (!item) {
    return (
      <View style={[styles.screen, { paddingTop: insets.top }]}>
        <ActivityIndicator color={colors.fgSecondary} style={{ flex: 1 }} />
      </View>
    );
  }

  const typeLabel = TYPES[item.type_id] ?? 'Other';
  const sub = [item.brand, item.model].filter(Boolean).join(' · ');

  return (
    <View style={styles.screen}>
      <ScrollView
        contentContainerStyle={[styles.scroll, { paddingTop: insets.top + 8 }]}
        showsVerticalScrollIndicator={false}
      >
        {/* Nav row */}
        <View style={styles.navRow}>
          <Pressable
            style={({ pressed }) => [styles.backBtn, pressed && { opacity: 0.7 }]}
            onPress={() => navigation.goBack()}
          >
            <Feather name="chevron-left" size={18} color={colors.fgPrimary} />
          </Pressable>
          <Pressable
            style={({ pressed }) => [styles.editBtn, pressed && { opacity: 0.7 }]}
            onPress={() => setEditOpen(true)}
          >
            <Text style={styles.editBtnLabel}>Edit</Text>
          </Pressable>
        </View>

        {/* Hero */}
        <View style={styles.hero}>
          <View style={styles.heroBubble}>
            <GearIcon typeId={item.type_id} size={36} color={palette.cream100} />
          </View>
          <View style={styles.typeBadge}>
            <Text style={styles.typeBadgeText}>{typeLabel}</Text>
          </View>
          <Text style={styles.heroName}>{item.name}</Text>
          {sub ? <Text style={styles.heroSub}>{sub}</Text> : null}
          {item.year ? (
            <Text style={styles.heroYear}>Acquired {item.year}</Text>
          ) : null}
        </View>

        {/* Notes card */}
        {item.notes ? (
          <View style={styles.card}>
            <Text style={styles.cardLabel}>NOTES</Text>
            <Text style={styles.cardBody}>{item.notes}</Text>
          </View>
        ) : null}

        {/* Extractions placeholder */}
        <View style={[styles.card, styles.placeholderCard]}>
          <Text style={styles.cardLabel}>EXTRACTIONS</Text>
          <Text style={styles.placeholderEmoji}>☕</Text>
          <Text style={styles.placeholderTitle}>No shots logged yet</Text>
          <Text style={styles.placeholderSub}>Shots using this gear will appear here.</Text>
        </View>

        {/* Remove button */}
        <Pressable
          style={({ pressed }) => [
            styles.removeBtn,
            pressed && { opacity: 0.75 },
            deleting && { opacity: 0.38 },
          ]}
          onPress={handleDelete}
          disabled={deleting}
        >
          {deleting ? (
            <ActivityIndicator color={palette.error500} />
          ) : (
            <Text style={styles.removeBtnLabel}>Remove from my gear</Text>
          )}
        </Pressable>
      </ScrollView>

      {editOpen && (
        <GearSheet
          editItem={item}
          onClose={() => setEditOpen(false)}
          onSaved={updated => {
            updateGear(updated);
            setEditOpen(false);
            toast.show(`${updated.name} updated ✓`);
          }}
        />
      )}

      <Animated.View style={[styles.toast, toast.animatedStyle]} pointerEvents="none">
        <Text style={styles.toastText}>{toast.msg}</Text>
      </Animated.View>
    </View>
  );
}

const styles = StyleSheet.create({
  screen: { flex: 1, backgroundColor: colors.bgApp },
  scroll: { paddingBottom: 48, paddingHorizontal: spacing[5] },

  navRow: {
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'space-between',
    marginBottom: 24,
  },
  backBtn: {
    width: 38,
    height: 38,
    borderRadius: 19,
    backgroundColor: palette.cream300,
    alignItems: 'center',
    justifyContent: 'center',
  },
  editBtn: {
    height: 34,
    paddingHorizontal: 16,
    borderRadius: radii.full,
    backgroundColor: palette.cream300,
    alignItems: 'center',
    justifyContent: 'center',
  },
  editBtnLabel: {
    fontFamily: 'DMSans_700Bold',
    fontSize: 13,
    color: colors.fgPrimary,
  },

  hero: { alignItems: 'center', paddingBottom: 24 },
  heroBubble: {
    width: 72,
    height: 72,
    borderRadius: 27,
    backgroundColor: palette.espresso800,
    alignItems: 'center',
    justifyContent: 'center',
    marginBottom: 14,
    shadowColor: palette.espresso800,
    shadowOffset: { width: 0, height: 4 },
    shadowOpacity: 0.25,
    shadowRadius: 12,
    elevation: 6,
  },
  typeBadge: {
    backgroundColor: palette.cream300,
    paddingHorizontal: 10,
    paddingVertical: 4,
    borderRadius: radii.full,
    marginBottom: 10,
  },
  typeBadgeText: {
    fontFamily: 'DMSans_500Medium',
    fontSize: 11,
    color: palette.espresso700,
  },
  heroName: {
    fontFamily: 'PlayfairDisplay_700Bold',
    fontSize: 26,
    color: colors.fgPrimary,
    textAlign: 'center',
    letterSpacing: -0.3,
    marginBottom: 4,
  },
  heroSub: {
    fontFamily: 'DMSans_400Regular',
    fontSize: 14,
    color: colors.fgSecondary,
    marginBottom: 2,
  },
  heroYear: {
    fontFamily: 'DMSans_400Regular',
    fontSize: 13,
    color: colors.fgTertiary,
  },

  card: {
    backgroundColor: colors.bgCard,
    borderRadius: radii.lg,
    padding: 18,
    marginBottom: 12,
    shadowColor: palette.espresso800,
    shadowOffset: { width: 0, height: 2 },
    shadowOpacity: 0.07,
    shadowRadius: 8,
    elevation: 2,
  },
  cardLabel: {
    fontFamily: 'DMSans_700Bold',
    fontSize: 10,
    color: palette.espresso400,
    letterSpacing: 0.8,
    textTransform: 'uppercase',
    marginBottom: 10,
  },
  cardBody: {
    fontFamily: 'DMSans_400Regular',
    fontSize: 14,
    color: colors.fgPrimary,
    lineHeight: 14 * 1.65,
  },

  placeholderCard: { alignItems: 'center', paddingVertical: 24 },
  placeholderEmoji: { fontSize: 28, marginBottom: 8 },
  placeholderTitle: {
    fontFamily: 'DMSans_700Bold',
    fontSize: 14,
    color: colors.fgPrimary,
    marginBottom: 4,
  },
  placeholderSub: {
    fontFamily: 'DMSans_400Regular',
    fontSize: 12,
    color: colors.fgTertiary,
    textAlign: 'center',
  },

  removeBtn: {
    height: 48,
    borderRadius: radii.xl,
    backgroundColor: palette.error100,
    borderWidth: 1.5,
    borderColor: palette.error100,
    alignItems: 'center',
    justifyContent: 'center',
    marginTop: 8,
  },
  removeBtnLabel: {
    fontFamily: 'DMSans_700Bold',
    fontSize: 15,
    color: palette.error500,
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
