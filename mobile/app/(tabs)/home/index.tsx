import { useState, useRef } from 'react';
import {
  View,
  Text,
  ScrollView,
  StyleSheet,
  Pressable,
  Animated,
} from 'react-native';
import { SafeAreaView } from 'react-native-safe-area-context';
import { useSafeAreaInsets } from 'react-native-safe-area-context';
import { Feather } from '@expo/vector-icons';
import { palette, spacing, radii } from '@/src/theme';
import { useExtractions } from '@/src/context/ExtractionsContext';
import { Extraction, computeZone } from '@/src/api/extractions';
import { ExtractionModal } from './ExtractionModal';

// ─── Helpers ─────────────────────────────────────────────────────────────────

function greeting(): string {
  const h = new Date().getHours();
  if (h < 12) return 'Good morning';
  if (h < 17) return 'Good afternoon';
  return 'Good evening';
}

function timeAgo(iso: string): string {
  const diff = (Date.now() - new Date(iso).getTime()) / 1000;
  if (diff < 60) return 'just now';
  if (diff < 3600) return `${Math.floor(diff / 60)}m ago`;
  if (diff < 86400) return `${Math.floor(diff / 3600)}h ago`;
  return `${Math.floor(diff / 86400)}d ago`;
}

function formatRatio(doseIn: number, yieldOut: number): string {
  return `1:${(yieldOut / doseIn).toFixed(1)}`;
}

// ─── Zone badge ───────────────────────────────────────────────────────────────

function ZoneBadge({ time, target }: { time: number; target: number }) {
  const zone = computeZone(time, target);

  let bg: string = palette.matcha100;
  let fg: string = palette.matcha700;
  let label = '✓ On target';

  if (zone === 'over') {
    bg = palette.error100;
    fg = palette.error500;
    label = 'Over';
  } else if (zone === 'under') {
    bg = palette.caramel100;
    fg = palette.caramel600;
    label = 'Under';
  }

  return (
    <View style={[zoneBadgeStyles.badge, { backgroundColor: bg }]}>
      <Text style={[zoneBadgeStyles.text, { color: fg }]}>{label}</Text>
    </View>
  );
}

const zoneBadgeStyles = StyleSheet.create({
  badge: {
    borderRadius: radii.full,
    paddingVertical: 3,
    paddingHorizontal: 10,
  },
  text: {
    fontFamily: 'DMSans_700Bold',
    fontSize: 11,
  },
});

// ─── Recent extraction card ───────────────────────────────────────────────────

function RecentCard({ extraction }: { extraction: Extraction }) {
  const ratio = formatRatio(extraction.dose_in, extraction.yield_out);

  return (
    <View style={cardStyles.card}>
      <View style={cardStyles.topRow}>
        <View style={cardStyles.topLeft}>
          <Text style={cardStyles.beanName}>{extraction.bean.name}</Text>
          <Text style={cardStyles.meta}>
            {extraction.bean.roaster
              ? `${extraction.bean.roaster} · `
              : ''}
            {timeAgo(extraction.created_at)}
          </Text>
        </View>
        <ZoneBadge time={extraction.time} target={extraction.target_time} />
      </View>

      <View style={cardStyles.statsRow}>
        <StatCol label="DOSE" value={`${extraction.dose_in}`} unit="g" />
        <View style={cardStyles.statDivider} />
        <StatCol label="YIELD" value={`${extraction.yield_out}`} unit="g" />
        <View style={cardStyles.statDivider} />
        <StatCol label="TIME" value={`${extraction.time.toFixed(1)}`} unit="s" />
        <View style={cardStyles.statDivider} />
        <StatCol label="RATIO" value={ratio} unit="" accent />
      </View>

      {extraction.tasting_note ? (
        <Text style={cardStyles.note}>"{extraction.tasting_note}"</Text>
      ) : null}
    </View>
  );
}

function StatCol({
  label,
  value,
  unit,
  accent,
}: {
  label: string;
  value: string;
  unit: string;
  accent?: boolean;
}) {
  return (
    <View style={cardStyles.statCol}>
      <View style={cardStyles.statValRow}>
        <Text style={[cardStyles.statVal, accent && { color: palette.caramel500 }]}>
          {value}
        </Text>
        {!!unit && <Text style={cardStyles.statUnit}>{unit}</Text>}
      </View>
      <Text style={cardStyles.statLabel}>{label}</Text>
    </View>
  );
}

const cardStyles = StyleSheet.create({
  card: {
    backgroundColor: palette.cream200,
    borderRadius: 20,
    padding: 18,
    paddingBottom: 16,
    shadowColor: '#1C0F07',
    shadowOffset: { width: 0, height: 2 },
    shadowOpacity: 0.06,
    shadowRadius: 5,
    elevation: 2,
  },
  topRow: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'flex-start',
    marginBottom: 14,
  },
  topLeft: {
    flex: 1,
    marginRight: 10,
  },
  beanName: {
    fontFamily: 'DMSans_700Bold',
    fontSize: 14,
    color: palette.espresso800,
  },
  meta: {
    fontFamily: 'DMSans_400Regular',
    fontSize: 12,
    color: palette.espresso400,
    marginTop: 2,
  },
  statsRow: {
    flexDirection: 'row',
    marginBottom: 12,
  },
  statDivider: {
    width: 1,
    backgroundColor: palette.cream400,
    marginVertical: 2,
  },
  statCol: {
    flex: 1,
    alignItems: 'center',
  },
  statValRow: {
    flexDirection: 'row',
    alignItems: 'flex-end',
    gap: 1,
  },
  statVal: {
    fontFamily: 'JetBrainsMono_600SemiBold',
    fontSize: 18,
    letterSpacing: -0.5,
    color: palette.espresso800,
  },
  statUnit: {
    fontFamily: 'JetBrainsMono_400Regular',
    fontSize: 10,
    color: palette.espresso400,
    marginBottom: 2,
  },
  statLabel: {
    fontFamily: 'DMSans_700Bold',
    fontSize: 9,
    letterSpacing: 0.07 * 9,
    textTransform: 'uppercase',
    color: palette.espresso400,
    marginTop: 4,
  },
  note: {
    fontFamily: 'DMSans_400Regular',
    fontSize: 12,
    fontStyle: 'italic',
    color: palette.espresso400,
    lineHeight: 1.45 * 12,
  },
});

// ─── Home screen ──────────────────────────────────────────────────────────────

export default function HomeScreen() {
  const insets = useSafeAreaInsets();
  const { extractions } = useExtractions();
  const [modalOpen, setModalOpen] = useState(false);

  const lastExtraction = extractions[0] ?? null;
  const recent = extractions.slice(0, 3);

  const lastSubline = lastExtraction
    ? `Last: ${lastExtraction.time.toFixed(1)}s · ${formatRatio(lastExtraction.dose_in, lastExtraction.yield_out)}`
    : null;

  const beanPillName = lastExtraction?.bean.name ?? 'Select beans';

  const cardPressAnim = useRef(new Animated.Value(1)).current;

  function onCardPressIn() {
    Animated.timing(cardPressAnim, {
      toValue: 0.98,
      duration: 120,
      useNativeDriver: true,
    }).start();
  }

  function onCardPressOut() {
    Animated.timing(cardPressAnim, {
      toValue: 1,
      duration: 120,
      useNativeDriver: true,
    }).start();
  }

  return (
    <View style={[styles.root, { backgroundColor: palette.cream100 }]}>
      <ScrollView
        style={styles.scroll}
        contentContainerStyle={{ paddingBottom: 32 }}
        showsVerticalScrollIndicator={false}
      >
        {/* Header */}
        <View
          style={[
            styles.header,
            { paddingTop: insets.top + 12 },
          ]}
        >
          <View style={styles.headerText}>
            <Text style={styles.greetingLabel}>{greeting()}</Text>
            <Text style={styles.headlineReady}>Ready to pull</Text>
            <Text style={styles.headlineShot}>a shot?</Text>
          </View>
          {/* Bean mark logo placeholder */}
          <View style={styles.logoMark}>
            <View style={styles.logoOuter}>
              <View style={styles.logoInner} />
            </View>
          </View>
        </View>

        {/* Extraction invitation card */}
        <Pressable
          onPressIn={onCardPressIn}
          onPressOut={onCardPressOut}
          onPress={() => setModalOpen(true)}
        >
          <Animated.View
            style={[styles.inviteCard, { transform: [{ scale: cardPressAnim }] }]}
          >
            {/* Decorative rings */}
            <View style={styles.decRing1} pointerEvents="none" />
            <View style={styles.decRing2} pointerEvents="none" />

            {/* Bean pill */}
            <View style={styles.beanPill}>
              <View style={styles.beanTear} />
              <Text style={styles.beanPillText} numberOfLines={1}>
                {beanPillName}
              </Text>
            </View>

            <Text style={styles.inviteHeadline}>Pull a shot.</Text>

            {lastSubline && (
              <Text style={styles.inviteSubline}>{lastSubline}</Text>
            )}

            <View style={styles.startBtn}>
              <Feather name="play" size={18} color="#fff" />
              <Text style={styles.startBtnText}>Start extraction</Text>
            </View>
          </Animated.View>
        </Pressable>

        {/* Recent extractions */}
        {recent.length > 0 && (
          <View style={styles.recentSection}>
            <View style={styles.recentHeader}>
              <Text style={styles.recentTitle}>Recent extractions</Text>
              <Pressable style={styles.allLink}>
                <Text style={styles.allLinkText}>All</Text>
                <Feather name="chevron-right" size={14} color={palette.matcha500} />
              </Pressable>
            </View>
            <View style={styles.cardList}>
              {recent.map(e => (
                <RecentCard key={e.id} extraction={e} />
              ))}
            </View>
          </View>
        )}
      </ScrollView>

      <ExtractionModal
        visible={modalOpen}
        onClose={() => setModalOpen(false)}
        lastExtraction={lastExtraction}
      />
    </View>
  );
}

// ─── Styles ───────────────────────────────────────────────────────────────────

const styles = StyleSheet.create({
  root: {
    flex: 1,
  },
  scroll: {
    flex: 1,
  },

  // Header
  header: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'flex-start',
    paddingHorizontal: 24,
    paddingBottom: 24,
  },
  headerText: {
    flex: 1,
  },
  greetingLabel: {
    fontFamily: 'DMSans_500Medium',
    fontSize: 13,
    color: palette.espresso400,
    marginBottom: 4,
  },
  headlineReady: {
    fontFamily: 'PlayfairDisplay_900Black',
    fontSize: 32,
    letterSpacing: -0.8,
    lineHeight: 32 * 1.1,
    color: palette.espresso800,
  },
  headlineShot: {
    fontFamily: 'PlayfairDisplay_700Bold_Italic',
    fontSize: 32,
    letterSpacing: -0.8,
    lineHeight: 32 * 1.1,
    color: palette.caramel500,
  },

  // Logo mark (simplified SVG stand-in)
  logoMark: {
    marginTop: 4,
  },
  logoOuter: {
    width: 34,
    height: 44,
    borderRadius: 17,
    backgroundColor: palette.espresso800,
    alignItems: 'center',
    justifyContent: 'center',
  },
  logoInner: {
    width: 14,
    height: 18,
    borderRadius: 7,
    backgroundColor: palette.caramel400,
  },

  // Invite card
  inviteCard: {
    marginHorizontal: 18,
    backgroundColor: palette.espresso800,
    borderRadius: 28,
    padding: 22,
    overflow: 'hidden',
    shadowColor: '#1C0F07',
    shadowOffset: { width: 0, height: 8 },
    shadowOpacity: 0.20,
    shadowRadius: 20,
    elevation: 8,
  },
  decRing1: {
    position: 'absolute',
    right: -40,
    top: -40,
    width: 180,
    height: 180,
    borderRadius: 90,
    borderWidth: 18,
    borderColor: 'rgba(255,255,255,0.03)',
  },
  decRing2: {
    position: 'absolute',
    right: 16,
    top: 16,
    width: 100,
    height: 100,
    borderRadius: 50,
    borderWidth: 16,
    borderColor: 'rgba(255,255,255,0.03)',
  },
  beanPill: {
    flexDirection: 'row',
    alignItems: 'center',
    alignSelf: 'flex-start',
    gap: 6,
    backgroundColor: palette.espresso700,
    borderRadius: radii.full,
    paddingVertical: 5,
    paddingLeft: 8,
    paddingRight: 12,
    marginBottom: 18,
  },
  beanTear: {
    width: 12,
    height: 15,
    borderTopLeftRadius: 6,
    borderTopRightRadius: 6,
    borderBottomLeftRadius: 4,
    borderBottomRightRadius: 4,
    backgroundColor: palette.caramel400,
  },
  beanPillText: {
    fontFamily: 'DMSans_700Bold',
    fontSize: 12,
    color: palette.cream500,
    maxWidth: 200,
  },
  inviteHeadline: {
    fontFamily: 'PlayfairDisplay_700Bold',
    fontSize: 28,
    letterSpacing: -0.4,
    lineHeight: 28 * 1.2,
    color: palette.cream100,
    marginBottom: 6,
  },
  inviteSubline: {
    fontFamily: 'DMSans_400Regular',
    fontSize: 13,
    color: palette.espresso400,
    marginBottom: 20,
  },
  startBtn: {
    height: 50,
    borderRadius: radii.full,
    backgroundColor: palette.matcha500,
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'center',
    gap: 8,
    shadowColor: palette.matcha500,
    shadowOffset: { width: 0, height: 4 },
    shadowOpacity: 0.40,
    shadowRadius: 8,
    elevation: 4,
  },
  startBtnText: {
    fontFamily: 'DMSans_700Bold',
    fontSize: 15,
    color: '#fff',
  },

  // Recent
  recentSection: {
    paddingHorizontal: 20,
    paddingTop: 32,
    paddingBottom: 0,
  },
  recentHeader: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    marginBottom: 12,
  },
  recentTitle: {
    fontFamily: 'PlayfairDisplay_700Bold',
    fontSize: 20,
    color: palette.espresso800,
  },
  allLink: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: 2,
  },
  allLinkText: {
    fontFamily: 'DMSans_700Bold',
    fontSize: 13,
    color: palette.matcha500,
  },
  cardList: {
    gap: 10,
  },
});
