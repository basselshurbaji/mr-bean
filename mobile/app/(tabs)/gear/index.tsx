import {
  Animated,
  FlatList,
  Pressable,
  StyleSheet,
  Text,
  View,
} from 'react-native';
import { useRef, useState } from 'react';
import { useRouter } from 'expo-router';
import { useSafeAreaInsets } from 'react-native-safe-area-context';
import { colors, palette, radii, spacing } from '@/src/theme';
import { useGear } from '@/src/context/GearContext';
import { GearItem, Station } from '@/src/api/gear';
import GearIcon from '@/src/components/GearIcon';
import GearSheet from './GearSheet';
import StationSheet from './StationSheet';

const TYPES = [
  { id: 'machine',     label: 'Espresso machine'  },
  { id: 'grinder',     label: 'Grinder'           },
  { id: 'scale',       label: 'Scale'             },
  { id: 'portafilter', label: 'Portafilter'       },
  { id: 'tamper',      label: 'Tamper'            },
  { id: 'distributor', label: 'Distribution tool' },
  { id: 'wdt',         label: 'WDT tool'          },
  { id: 'basket',      label: 'Basket'            },
  { id: 'puckscreen',  label: 'Puck screen'       },
  { id: 'other',       label: 'Other'             },
];

const typeLabel = (id: string) => TYPES.find(t => t.id === id)?.label ?? 'Other';

// ─── Toast ────────────────────────────────────────────────────────────────────

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

// ─── Gear card ────────────────────────────────────────────────────────────────

function GearCard({ item, onPress }: { item: GearItem; onPress: () => void }) {
  const scale = useRef(new Animated.Value(1)).current;
  const sub = [item.brand, item.model].filter(Boolean).join(' · ');

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
        <View style={styles.cardBubble}>
          <GearIcon typeId={item.type_id} size={25} color={palette.espresso800} />
        </View>
        <View style={styles.cardText}>
          <Text style={styles.cardName} numberOfLines={1}>{item.name}</Text>
          {sub ? (
            <Text style={styles.cardSub} numberOfLines={1}>{sub}</Text>
          ) : null}
        </View>
        <View style={styles.cardMeta}>
          <View style={styles.typeBadge}>
            <Text style={styles.typeBadgeText}>{typeLabel(item.type_id)}</Text>
          </View>
          {item.year ? (
            <Text style={styles.cardYear}>{item.year}</Text>
          ) : null}
        </View>
      </Animated.View>
    </Pressable>
  );
}

// ─── Station card ─────────────────────────────────────────────────────────────

function StationCard({ station, onPress }: { station: Station; onPress: () => void }) {
  const scale = useRef(new Animated.Value(1)).current;
  const visibleGear = station.gear.slice(0, 7);
  const overflow = station.gear.length - 7;

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
      <Animated.View style={[styles.card, styles.stationCard, { transform: [{ scale }] }]}>
        <View style={styles.stationRow}>
          <View style={styles.stationTextWrap}>
            <Text style={styles.stationName} numberOfLines={1}>{station.name}</Text>
            <Text style={styles.stationCount}>
              {station.gear.length} item{station.gear.length !== 1 ? 's' : ''}
            </Text>
          </View>
          <Text style={styles.chevron}>›</Text>
        </View>
        {station.gear.length > 0 && (
          <View style={styles.iconStrip}>
            {visibleGear.map(g => (
              <View key={g.id} style={styles.iconTile}>
                <GearIcon typeId={g.type_id} size={18} color={palette.espresso800} />
              </View>
            ))}
            {overflow > 0 && (
              <View style={styles.iconTile}>
                <Text style={styles.iconOverflow}>+{overflow}</Text>
              </View>
            )}
          </View>
        )}
      </Animated.View>
    </Pressable>
  );
}

// ─── Screen ───────────────────────────────────────────────────────────────────

type ActiveSheet =
  | { type: 'add-gear' }
  | { type: 'add-station' }
  | { type: 'edit-station'; station: Station };

export default function GearScreen() {
  const insets = useSafeAreaInsets();
  const router = useRouter();
  const { gear, stations, addGear, addStation, updateStation, removeStation } = useGear();
  const toast = useToast();

  const [activeTab, setActiveTab] = useState<'gear' | 'stations'>('gear');
  const [filterType, setFilterType] = useState('all');
  const [activeSheet, setActiveSheet] = useState<ActiveSheet | null>(null);

  const presentTypes = ['all', ...Array.from(new Set(gear.map(g => g.type_id)))];
  const filteredGear =
    filterType === 'all' ? gear : gear.filter(g => g.type_id === filterType);

  function closeSheet() {
    setActiveSheet(null);
  }

  return (
    <View style={styles.screen}>
      {/* Header */}
      <View style={[styles.header, { paddingTop: insets.top + 12 }]}>
        <Text style={styles.title}>My Gear</Text>
        <Text style={styles.subtitle}>
          {gear.length} piece{gear.length !== 1 ? 's' : ''} · {stations.length} station{stations.length !== 1 ? 's' : ''}
        </Text>
      </View>

      {/* Segment control */}
      <View style={styles.segmentWrap}>
        <View style={styles.segment}>
          {(['gear', 'stations'] as const).map(tab => (
            <Pressable
              key={tab}
              style={[styles.segmentBtn, activeTab === tab && styles.segmentBtnActive]}
              onPress={() => setActiveTab(tab)}
            >
              <Text style={[styles.segmentLabel, activeTab === tab && styles.segmentLabelActive]}>
                {tab === 'gear' ? 'Gear' : 'Stations'}
              </Text>
            </Pressable>
          ))}
        </View>
      </View>

      {activeTab === 'gear' ? (
        <>
          {/* Filter chips */}
          <FlatList
            horizontal
            data={presentTypes}
            keyExtractor={t => t}
            showsHorizontalScrollIndicator={false}
            style={styles.chipsRow}
            contentContainerStyle={styles.chips}
            renderItem={({ item: t }) => (
              <Pressable
                style={[styles.chip, filterType === t && styles.chipActive]}
                onPress={() => setFilterType(t)}
              >
                <Text style={[styles.chipLabel, filterType === t && styles.chipLabelActive]}>
                  {t === 'all' ? 'All' : typeLabel(t)}
                </Text>
              </Pressable>
            )}
          />

          {/* Gear list */}
          <FlatList
            data={filteredGear}
            keyExtractor={g => g.id}
            contentContainerStyle={styles.list}
            showsVerticalScrollIndicator={false}
            ListEmptyComponent={
              <View style={styles.empty}>
                <Text style={styles.emptyEmoji}>⚙️</Text>
                <Text style={styles.emptyTitle}>
                  {filterType === 'all' ? 'No gear yet' : 'No matches'}
                </Text>
                <Text style={styles.emptySub}>
                  {filterType === 'all'
                    ? 'Add your first piece of equipment below.'
                    : 'Try a different filter or add new gear.'}
                </Text>
              </View>
            }
            ListFooterComponent={
              <Pressable
                style={({ pressed }) => [styles.addGearBtn, pressed && { backgroundColor: palette.cream300 }]}
                onPress={() => setActiveSheet({ type: 'add-gear' })}
              >
                <Text style={styles.addGearLabel}>+ Add gear</Text>
              </Pressable>
            }
            renderItem={({ item }) => (
              <GearCard
                item={item}
                // eslint-disable-next-line @typescript-eslint/no-explicit-any
                onPress={() => (router.push as any)({ pathname: '/(tabs)/gear/[id]', params: { id: item.id } })}
              />
            )}
          />
        </>
      ) : (
        /* Stations tab */
        <FlatList
          data={stations}
          keyExtractor={s => s.id}
          contentContainerStyle={styles.list}
          showsVerticalScrollIndicator={false}
          ListEmptyComponent={
            <View style={styles.empty}>
              <Text style={styles.emptyEmoji}>🗂️</Text>
              <Text style={styles.emptyTitle}>No stations yet</Text>
              <Text style={styles.emptySub}>
                Create a station to group your gear into a quick-select preset.
              </Text>
            </View>
          }
          ListFooterComponent={
            <Pressable
              style={({ pressed }) => [styles.addStationBtn, pressed && { backgroundColor: palette.cream300 }]}
              onPress={() => setActiveSheet({ type: 'add-station' })}
            >
              <Text style={styles.addStationLabel}>+ New station</Text>
            </Pressable>
          }
          renderItem={({ item }) => (
            <StationCard
              station={item}
              onPress={() => setActiveSheet({ type: 'edit-station', station: item })}
            />
          )}
        />
      )}

      {/* Sheets */}
      {activeSheet?.type === 'add-gear' && (
        <GearSheet
          onClose={closeSheet}
          onSaved={item => {
            addGear(item);
            closeSheet();
            toast.show(`${item.name} added ✓`);
          }}
        />
      )}

      {activeSheet?.type === 'add-station' && (
        <StationSheet
          gear={gear}
          onClose={closeSheet}
          onSaved={station => {
            addStation(station);
            closeSheet();
            toast.show(`${station.name} created ✓`);
          }}
        />
      )}

      {activeSheet?.type === 'edit-station' && (
        <StationSheet
          gear={gear}
          editStation={activeSheet.station}
          onClose={closeSheet}
          onSaved={station => {
            updateStation(station);
            closeSheet();
            toast.show(`${station.name} updated ✓`);
          }}
          onDeleted={id => {
            removeStation(id);
            closeSheet();
            toast.show('Station deleted');
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
    color: colors.fgTertiary,
    marginTop: 2,
  },

  segmentWrap: {
    paddingHorizontal: spacing[5],
    paddingBottom: spacing[3],
  },
  segment: {
    flexDirection: 'row',
    backgroundColor: palette.cream300,
    borderRadius: 14,
    padding: 4,
  },
  segmentBtn: {
    flex: 1,
    height: 36,
    borderRadius: 10,
    alignItems: 'center',
    justifyContent: 'center',
  },
  segmentBtnActive: {
    backgroundColor: palette.cream100,
    shadowColor: palette.espresso800,
    shadowOffset: { width: 0, height: 1 },
    shadowOpacity: 0.1,
    shadowRadius: 4,
    elevation: 2,
  },
  segmentLabel: {
    fontFamily: 'DMSans_700Bold',
    fontSize: 13,
    color: palette.espresso400,
  },
  segmentLabelActive: { color: palette.espresso800 },

  chipsRow: {
    flexGrow: 0,
  },
  chips: {
    paddingHorizontal: spacing[5],
    paddingBottom: spacing[4],
    gap: 8,
  },
  chip: {
    height: 32,
    paddingHorizontal: 14,
    borderRadius: radii.full,
    borderWidth: 1.5,
    borderColor: palette.cream400,
    alignItems: 'center',
    justifyContent: 'center',
  },
  chipActive: {
    backgroundColor: palette.espresso800,
    borderColor: palette.espresso800,
  },
  chipLabel: {
    fontFamily: 'DMSans_700Bold',
    fontSize: 12,
    color: palette.espresso500,
  },
  chipLabelActive: { color: palette.cream100 },

  list: {
    paddingHorizontal: spacing[5],
    paddingBottom: 120,
    gap: 10,
  },

  card: {
    backgroundColor: palette.cream200,
    borderRadius: radii.lg,
    padding: 14,
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
  cardBubble: {
    width: 50,
    height: 50,
    borderRadius: 19,
    backgroundColor: palette.cream300,
    alignItems: 'center',
    justifyContent: 'center',
  },
  cardText: { flex: 1 },
  cardName: {
    fontFamily: 'DMSans_700Bold',
    fontSize: 15,
    color: colors.fgPrimary,
  },
  cardSub: {
    fontFamily: 'DMSans_400Regular',
    fontSize: 12,
    color: colors.fgSecondary,
    marginTop: 2,
  },
  cardMeta: { alignItems: 'flex-end', gap: 4 },
  typeBadge: {
    backgroundColor: palette.cream300,
    paddingHorizontal: 8,
    paddingVertical: 3,
    borderRadius: radii.full,
  },
  typeBadgeText: {
    fontFamily: 'DMSans_500Medium',
    fontSize: 10,
    color: palette.espresso700,
  },
  cardYear: {
    fontFamily: 'DMSans_400Regular',
    fontSize: 11,
    color: palette.cream600,
  },

  stationCard: {
    flexDirection: 'column',
    padding: 18,
    paddingBottom: 16,
    gap: 10,
  },
  stationRow: {
    flexDirection: 'row',
    alignItems: 'center',
  },
  stationTextWrap: { flex: 1 },
  stationName: {
    fontFamily: 'DMSans_700Bold',
    fontSize: 15,
    color: colors.fgPrimary,
  },
  stationCount: {
    fontFamily: 'DMSans_400Regular',
    fontSize: 12,
    color: colors.fgSecondary,
    marginTop: 1,
  },
  chevron: {
    fontSize: 20,
    color: colors.fgTertiary,
    lineHeight: 24,
  },
  iconStrip: {
    flexDirection: 'row',
    gap: 6,
    flexWrap: 'wrap',
  },
  iconTile: {
    width: 36,
    height: 36,
    borderRadius: 10,
    backgroundColor: palette.cream300,
    alignItems: 'center',
    justifyContent: 'center',
  },
  iconOverflow: {
    fontFamily: 'DMSans_700Bold',
    fontSize: 10,
    color: colors.fgSecondary,
  },

  addStationBtn: {
    height: 54,
    borderRadius: radii.lg,
    borderWidth: 1.5,
    borderColor: palette.cream500,
    borderStyle: 'dashed',
    alignItems: 'center',
    justifyContent: 'center',
    marginTop: 4,
  },
  addStationLabel: {
    fontFamily: 'DMSans_700Bold',
    fontSize: 14,
    color: colors.fgSecondary,
  },

  addGearBtn: {
    height: 54,
    borderRadius: radii.lg,
    borderWidth: 1.5,
    borderColor: palette.cream500,
    borderStyle: 'dashed',
    alignItems: 'center',
    justifyContent: 'center',
    marginTop: 4,
  },
  addGearLabel: {
    fontFamily: 'DMSans_700Bold',
    fontSize: 14,
    color: colors.fgSecondary,
  },

  empty: {
    alignItems: 'center',
    paddingTop: 52,
    paddingHorizontal: spacing[6],
  },
  emptyEmoji: { fontSize: 40, marginBottom: 12 },
  emptyTitle: {
    fontFamily: 'DMSans_700Bold',
    fontSize: 16,
    color: colors.fgPrimary,
    marginBottom: 6,
  },
  emptySub: {
    fontFamily: 'DMSans_400Regular',
    fontSize: 13,
    color: colors.fgTertiary,
    textAlign: 'center',
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
