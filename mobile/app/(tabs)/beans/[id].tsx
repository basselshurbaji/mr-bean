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
import { BeansStackParamList } from '@/src/navigation';
import { useRef, useState } from 'react';
import { useSafeAreaInsets } from 'react-native-safe-area-context';
import { Feather } from '@expo/vector-icons';
import { colors, palette, radii, spacing } from '@/src/theme';
import { useBeans } from '@/src/context/BeansContext';
import { beansApi, processLabel, roastLabel, roastColor } from '@/src/api/beans';
import RoastBubble from '@/src/components/RoastBubble';
import BeanSheet from './BeanSheet';

function useToast() {
  const opacity    = useRef(new Animated.Value(0)).current;
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

export default function BeanDetailScreen() {
  const route = useRoute<RouteProp<BeansStackParamList, 'BeanDetail'>>();
  const { id } = route.params;
  const navigation = useNavigation<NativeStackNavigationProp<BeansStackParamList>>();
  const insets  = useSafeAreaInsets();
  const { beans, updateBean, removeBean } = useBeans();
  const toast   = useToast();

  const bean = beans.find(b => b.id === id);

  const [editOpen,  setEditOpen]  = useState(false);
  const [deleting,  setDeleting]  = useState(false);

  async function handleDelete() {
    if (deleting || !bean) return;
    setDeleting(true);
    try {
      await beansApi.delete(bean.id);
      removeBean(bean.id);
      navigation.goBack();
    } catch {
      setDeleting(false);
      toast.show('Failed to remove bean.');
    }
  }

  if (!bean) {
    return (
      <View style={[styles.screen, { paddingTop: insets.top }]}>
        <ActivityIndicator color={colors.fgSecondary} style={{ flex: 1 }} />
      </View>
    );
  }

  const proc  = processLabel(bean.process);
  const roast = roastLabel(bean.roast_level);
  const rColor = roastColor(bean.roast_level);

  const tastingChips = bean.tasting_notes
    ? bean.tasting_notes.replace(/\.\s*$/, '').split(/,\s*/).filter(Boolean)
    : [];

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
          <RoastBubble roastId={bean.roast_level} size={72} />

          {/* Badge row */}
          {(proc || roast) ? (
            <View style={styles.badgeRow}>
              {proc ? (
                <View style={styles.procBadge}>
                  <Text style={styles.procBadgeText}>{proc}</Text>
                </View>
              ) : null}
              {roast ? (
                <View style={[styles.roastBadge, { backgroundColor: rColor }]}>
                  <Text style={styles.roastBadgeText}>{roast}</Text>
                </View>
              ) : null}
            </View>
          ) : null}

          <Text style={styles.heroName}>{bean.name}</Text>
          {bean.roaster ? <Text style={styles.heroRoaster}>{bean.roaster}</Text> : null}
          {bean.origin  ? <Text style={styles.heroOrigin}>{bean.origin}</Text>   : null}
        </View>

        {/* Tasting notes card */}
        {tastingChips.length > 0 ? (
          <View style={styles.card}>
            <Text style={styles.cardLabel}>TASTING NOTES</Text>
            <View style={styles.chipRow}>
              {tastingChips.map((chip, i) => (
                <View key={i} style={styles.tastingChip}>
                  <Text style={styles.tastingChipText}>{chip.trim()}</Text>
                </View>
              ))}
            </View>
          </View>
        ) : null}

        {/* Personal notes card */}
        {bean.notes ? (
          <View style={styles.card}>
            <Text style={styles.cardLabel}>NOTES</Text>
            <Text style={styles.cardBody}>{bean.notes}</Text>
          </View>
        ) : null}

        {/* Remove */}
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
            <Text style={styles.removeBtnLabel}>Remove this bean</Text>
          )}
        </Pressable>
      </ScrollView>

      {editOpen && (
        <BeanSheet
          editBean={bean}
          onClose={() => setEditOpen(false)}
          onSaved={updated => {
            updateBean(updated);
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
  scroll: { paddingHorizontal: spacing[6], paddingBottom: 48 },

  navRow: {
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'space-between',
    marginBottom: 20,
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

  hero: {
    alignItems: 'center',
    paddingBottom: 24,
    gap: 10,
  },
  badgeRow: {
    flexDirection: 'row',
    gap: 6,
    alignItems: 'center',
  },
  procBadge: {
    backgroundColor: palette.cream300,
    paddingHorizontal: 10,
    paddingVertical: 4,
    borderRadius: radii.full,
  },
  procBadgeText: {
    fontFamily: 'DMSans_700Bold',
    fontSize: 10,
    color: palette.espresso600,
  },
  roastBadge: {
    paddingHorizontal: 10,
    paddingVertical: 4,
    borderRadius: radii.full,
  },
  roastBadgeText: {
    fontFamily: 'DMSans_700Bold',
    fontSize: 10,
    color: palette.cream100,
  },
  heroName: {
    fontFamily: 'PlayfairDisplay_700Bold',
    fontSize: 26,
    color: colors.fgPrimary,
    textAlign: 'center',
    letterSpacing: -0.3,
    lineHeight: 26 * 1.15,
  },
  heroRoaster: {
    fontFamily: 'DMSans_400Regular',
    fontSize: 14,
    color: palette.espresso500,
    marginTop: -4,
  },
  heroOrigin: {
    fontFamily: 'DMSans_400Regular',
    fontSize: 13,
    color: palette.espresso300,
    marginTop: -4,
  },

  card: {
    backgroundColor: palette.cream200,
    borderRadius: radii.lg,
    padding: 18,
    paddingHorizontal: 20,
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
    marginBottom: 12,
  },
  cardBody: {
    fontFamily: 'DMSans_400Regular',
    fontSize: 14,
    color: colors.fgPrimary,
    lineHeight: 14 * 1.65,
  },

  chipRow: {
    flexDirection: 'row',
    flexWrap: 'wrap',
    gap: 6,
  },
  tastingChip: {
    backgroundColor: palette.cream300,
    paddingHorizontal: 10,
    paddingVertical: 4,
    borderRadius: radii.full,
  },
  tastingChipText: {
    fontFamily: 'DMSans_700Bold',
    fontSize: 12,
    color: palette.espresso600,
  },

  removeBtn: {
    height: 48,
    borderRadius: radii.xl,
    backgroundColor: palette.error100,
    borderWidth: 1.5,
    borderColor: palette.error100,
    alignItems: 'center',
    justifyContent: 'center',
    marginTop: 4,
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
