import React, { useState, useRef, useEffect, useCallback } from 'react';
import {
  View,
  Text,
  StyleSheet,
  Pressable,
  ScrollView,
  Modal,
  Animated,
  Dimensions,
  TextInput,
  Alert,
  RefreshControl,
} from 'react-native';
import { KeyboardAwareScrollView } from 'react-native-keyboard-controller';
import { Feather } from '@expo/vector-icons';
import Svg, { Circle } from 'react-native-svg';
import { useSafeAreaInsets } from 'react-native-safe-area-context';
import { palette, spacing, radii } from '@/src/theme';
import { Bean, roastColor, roastLabel } from '@/src/api/beans';
import { GearItem, Station } from '@/src/api/gear';
import { extractionApi, Extraction, computeZone } from '@/src/api/extractions';
import { useExtractions } from '@/src/context/ExtractionsContext';
import { useBeans } from '@/src/context/BeansContext';
import { useGear } from '@/src/context/GearContext';

// ─── Constants ───────────────────────────────────────────────────────────────

const { height: SCREEN_HEIGHT } = Dimensions.get('window');
const SHEET_HEIGHT = SCREEN_HEIGHT * 0.92;
const SUB_SHEET_MAX = SCREEN_HEIGHT * 0.88;

const RING_SIZE = 220;
const RING_CX = 110;
const RING_CY = 110;
const RING_R = 88;
const RING_SW = 10;
const CIRCUMFERENCE = 2 * Math.PI * RING_R;

// ─── Helpers ─────────────────────────────────────────────────────────────────

type TimerState = 'idle' | 'running' | 'done';

function fmtTime(s: number): string {
  const m = Math.floor(s / 60);
  const sec = Math.floor(s % 60);
  return `${String(m).padStart(2, '0')}:${String(sec).padStart(2, '0')}`;
}

function zoneColor(elapsed: number, target: number): string {
  const z = computeZone(elapsed, target);
  if (z === 'under') return palette.caramel400;
  if (z === 'perfect') return palette.matcha500;
  return palette.error500;
}

function zoneLabel(elapsed: number, target: number): string {
  const z = computeZone(elapsed, target);
  if (z === 'under') return 'Under';
  if (z === 'perfect') return 'On target';
  return 'Over';
}

function zoneBgColor(elapsed: number, target: number): string {
  const z = computeZone(elapsed, target);
  if (z === 'under') return palette.caramel100;
  if (z === 'perfect') return palette.matcha100;
  return palette.error100;
}

// ─── Sub-sheet: Beans ────────────────────────────────────────────────────────

interface BeanSheetProps {
  beans: Bean[];
  loading: boolean;
  selectedId: string | null;
  onSelect: (bean: Bean) => void;
  onClose: () => void;
  translateY: Animated.Value;
  onRefresh: () => void;
}

function BeanSheet({ beans, loading, selectedId, onSelect, onClose, translateY, onRefresh }: BeanSheetProps) {
  return (
    <Animated.View style={[styles.subSheet, { transform: [{ translateY }] }]}>
      <View style={styles.subHandle} />
      <View style={styles.subHeader}>
        <Text style={styles.subTitle}>Select beans</Text>
        <Pressable style={styles.subClose} onPress={onClose}>
          <Feather name="x" size={16} color={palette.espresso400} />
        </Pressable>
      </View>
      <ScrollView
        style={styles.subList}
        contentContainerStyle={styles.subListContent}
        refreshControl={
          <RefreshControl refreshing={loading} onRefresh={onRefresh} tintColor={palette.caramel400} />
        }
      >
        {loading && beans.length === 0 ? (
          <Text style={styles.emptyHint}>Loading beans…</Text>
        ) : beans.length === 0 ? (
          <Text style={styles.emptyHint}>No beans found. Add some in the Beans tab.</Text>
        ) : (
          beans.map(bean => {
            const selected = bean.id === selectedId;
            const rc = roastColor(bean.roast_level);
            const meta = [bean.roaster, bean.roast_level ? roastLabel(bean.roast_level) : null]
              .filter(Boolean)
              .join(' · ');
            return (
              <Pressable
                key={bean.id}
                style={[styles.beanRow, selected && styles.beanRowSel]}
                onPress={() => { onSelect(bean); onClose(); }}
              >
                <View style={[styles.beanIconBox, { backgroundColor: selected ? palette.espresso700 : palette.cream300 }]}>
                  <View style={[styles.beanTear, { backgroundColor: rc }]} />
                </View>
                <View style={styles.beanInfo}>
                  <Text style={[styles.beanRowName, { color: selected ? palette.cream100 : palette.espresso800 }]}>
                    {bean.name}
                  </Text>
                  {!!meta && (
                    <Text style={[styles.beanRowMeta, { color: selected ? palette.cream500 : palette.espresso400 }]}>
                      {meta}
                    </Text>
                  )}
                </View>
                {selected && <Feather name="check" size={16} color={palette.caramel400} />}
              </Pressable>
            );
          })
        )}
        <Pressable style={styles.addBeanBtn}>
          <Feather name="plus" size={14} color={palette.espresso400} />
          <Text style={styles.addBeanText}>Add new beans</Text>
        </Pressable>
      </ScrollView>
    </Animated.View>
  );
}

// ─── Sub-sheet: Gear ─────────────────────────────────────────────────────────

interface GearSheetProps {
  gear: GearItem[];
  stations: Station[];
  loading: boolean;
  selectedIds: string[];
  onSave: (ids: string[]) => void;
  onClose: () => void;
  translateY: Animated.Value;
  onRefresh: () => void;
}

function GearSheet({ gear, stations, loading, selectedIds, onSave, onClose, translateY, onRefresh }: GearSheetProps) {
  const [localIds, setLocalIds] = useState<string[]>(selectedIds);

  useEffect(() => {
    setLocalIds(selectedIds);
  }, [selectedIds]);

  function toggleGear(id: string) {
    setLocalIds(prev =>
      prev.includes(id) ? prev.filter(x => x !== id) : [...prev, id],
    );
  }

  function addStation(station: Station) {
    setLocalIds(prev => {
      const toAdd = station.gear.map(g => g.id).filter(id => !prev.includes(id));
      return [...prev, ...toAdd];
    });
  }

  return (
    <Animated.View style={[styles.subSheet, { transform: [{ translateY }] }]}>
      <View style={styles.subHandle} />
      <View style={styles.subHeader}>
        <Text style={styles.subTitle}>Gear used</Text>
        <Pressable style={styles.subClose} onPress={onClose}>
          <Feather name="x" size={16} color={palette.espresso400} />
        </Pressable>
      </View>

      <ScrollView
        style={styles.subList}
        contentContainerStyle={styles.subListContent}
        refreshControl={
          <RefreshControl refreshing={loading} onRefresh={onRefresh} tintColor={palette.caramel400} />
        }
      >
        {stations.length > 0 && (
          <>
            <Text style={styles.gearSectionLabel}>Stations</Text>
            <View style={styles.stationRow}>
              {stations.map(s => (
                <Pressable key={s.id} style={styles.stationBtn} onPress={() => addStation(s)}>
                  <Text style={styles.stationBtnText}>{s.name}</Text>
                </Pressable>
              ))}
            </View>
            <Text style={styles.gearSectionLabel}>Individual gear</Text>
          </>
        )}

        {loading && gear.length === 0 ? (
          <Text style={styles.emptyHint}>Loading gear…</Text>
        ) : gear.length === 0 ? (
          <Text style={styles.emptyHint}>No gear found. Add some in the My Gear tab.</Text>
        ) : (
          gear.map(item => {
            const selected = localIds.includes(item.id);
            return (
              <Pressable
                key={item.id}
                style={[styles.gearItemRow, selected && styles.gearItemRowSel]}
                onPress={() => toggleGear(item.id)}
              >
                <View style={styles.gearInfo}>
                  <Text style={[styles.gearItemName, { color: selected ? palette.cream100 : palette.espresso800 }]}>
                    {item.name}
                  </Text>
                  <Text style={[styles.gearItemType, { color: selected ? palette.cream500 : palette.espresso400 }]}>
                    {item.type_id}
                  </Text>
                </View>
                {selected && <Feather name="check" size={16} color={palette.caramel400} />}
              </Pressable>
            );
          })
        )}
      </ScrollView>

      <View style={styles.gearFooter}>
        <Pressable
          style={styles.gearSaveBtn}
          onPress={() => { onSave(localIds); onClose(); }}
        >
          <Text style={styles.gearSaveBtnText}>
            Done{localIds.length > 0 ? ` · ${localIds.length} item${localIds.length !== 1 ? 's' : ''} selected` : ''}
          </Text>
        </Pressable>
      </View>
    </Animated.View>
  );
}

// ─── Timer Ring ───────────────────────────────────────────────────────────────

interface RingProps {
  timerState: TimerState;
  elapsed: number;
  targetTime: number;
}

function TimerRing({ timerState, elapsed, targetTime }: RingProps) {
  const underFraction = Math.max(0, (targetTime - 4) / targetTime);
  const underLen = underFraction * CIRCUMFERENCE;
  const perfectLen = CIRCUMFERENCE - underLen;
  const underStartDeg = -90;
  const perfectStartDeg = -90 + underFraction * 360;

  const isActive = timerState === 'running' || timerState === 'done';
  const fraction = Math.min(elapsed / targetTime, 1.2);
  const progressLen = Math.min(fraction * CIRCUMFERENCE, CIRCUMFERENCE);
  const arcColor = isActive ? zoneColor(elapsed, targetTime) : palette.cream400;

  return (
    <Svg width={RING_SIZE} height={RING_SIZE}>
      {/* Track */}
      <Circle
        cx={RING_CX}
        cy={RING_CY}
        r={RING_R}
        stroke={palette.cream400}
        strokeWidth={RING_SW}
        fill="none"
      />
      {/* Under zone arc */}
      {underLen > 0 && (
        <Circle
          cx={RING_CX}
          cy={RING_CY}
          r={RING_R}
          stroke={palette.caramel300}
          strokeWidth={RING_SW}
          fill="none"
          strokeDasharray={[underLen, CIRCUMFERENCE - underLen]}
          rotation={underStartDeg}
          originX={RING_CX}
          originY={RING_CY}
          opacity={0.30}
        />
      )}
      {/* Perfect zone arc */}
      {perfectLen > 0 && (
        <Circle
          cx={RING_CX}
          cy={RING_CY}
          r={RING_R}
          stroke={palette.matcha400}
          strokeWidth={RING_SW}
          fill="none"
          strokeDasharray={[perfectLen, CIRCUMFERENCE - perfectLen]}
          rotation={perfectStartDeg}
          originX={RING_CX}
          originY={RING_CY}
          opacity={0.35}
        />
      )}
      {/* Progress arc */}
      {isActive && progressLen > 0 && (
        <Circle
          cx={RING_CX}
          cy={RING_CY}
          r={RING_R}
          stroke={arcColor}
          strokeWidth={RING_SW}
          fill="none"
          strokeDasharray={[progressLen, CIRCUMFERENCE - progressLen + 1]}
          strokeLinecap="round"
          rotation={-90}
          originX={RING_CX}
          originY={RING_CY}
        />
      )}
    </Svg>
  );
}

// ─── Main Modal ───────────────────────────────────────────────────────────────

export interface ExtractionModalProps {
  visible: boolean;
  onClose: () => void;
  lastExtraction?: Extraction | null;
}

export function ExtractionModal({
  visible,
  onClose,
  lastExtraction,
}: ExtractionModalProps) {
  const insets = useSafeAreaInsets();
  const { addExtraction } = useExtractions();
  const { beans, refresh: refreshBeans, loading: beansLoading } = useBeans();
  const { gear, stations, refresh: refreshGear, loading: gearLoading } = useGear();

  // ── Animations ──
  const sheetAnim = useRef(new Animated.Value(SHEET_HEIGHT)).current;
  const backdropAnim = useRef(new Animated.Value(0)).current;
  const beanSheetAnim = useRef(new Animated.Value(SUB_SHEET_MAX)).current;
  const gearSheetAnim = useRef(new Animated.Value(SUB_SHEET_MAX)).current;

  // ── Timer ──
  const [timerState, setTimerState] = useState<TimerState>('idle');
  const timerStateRef = useRef<TimerState>('idle');
  const [elapsed, setElapsed] = useState(0);
  const intervalRef = useRef<ReturnType<typeof setInterval> | null>(null);
  const runStartRef = useRef(0);

  // ── Form ──
  const [selectedBean, setSelectedBean] = useState<Bean | null>(null);
  const [doseIn, setDoseIn] = useState('');
  const [yieldOut, setYieldOut] = useState('');
  const [grindSize, setGrindSize] = useState('');
  const [targetTime, setTargetTime] = useState(27);
  const [preInfusion, setPreInfusion] = useState(false);
  const [selectedGearIds, setSelectedGearIds] = useState<string[]>([]);
  const [tastingNote, setTastingNote] = useState('');
  const [manualMode, setManualMode] = useState(false);
  const [manualSeconds, setManualSeconds] = useState('');

  // ── UI ──
  const [beanSheetOpen, setBeanSheetOpen] = useState(false);
  const [gearSheetOpen, setGearSheetOpen] = useState(false);
  const [saving, setSaving] = useState(false);

  // ── Prefill from last extraction ──
  useEffect(() => {
    if (visible && lastExtraction) {
      setTargetTime(Math.round(lastExtraction.target_time));
      setDoseIn(String(lastExtraction.dose_in));
      setGrindSize(lastExtraction.grind_size > 0 ? String(lastExtraction.grind_size) : '');
      setSelectedGearIds(lastExtraction.gear.map(g => g.id));
    }
  }, [visible, lastExtraction]);

  // ── Open / close — refresh data if empty when modal opens ──
  useEffect(() => {
    if (visible) {
      sheetAnim.setValue(SHEET_HEIGHT);
      backdropAnim.setValue(0);
      Animated.parallel([
        Animated.spring(sheetAnim, {
          toValue: 0,
          tension: 65,
          friction: 10,
          useNativeDriver: true,
        }),
        Animated.timing(backdropAnim, {
          toValue: 1,
          duration: 280,
          useNativeDriver: true,
        }),
      ]).start();
      if (beans.length === 0) refreshBeans();
      if (gear.length === 0) refreshGear();
    }
  }, [visible, sheetAnim, backdropAnim, beans.length, gear.length, refreshBeans, refreshGear]);

  function dismiss() {
    if (timerState === 'running') return;
    closeBeanSheet();
    closeGearSheet();
    Animated.parallel([
      Animated.timing(sheetAnim, {
        toValue: SHEET_HEIGHT,
        duration: 300,
        useNativeDriver: true,
      }),
      Animated.timing(backdropAnim, {
        toValue: 0,
        duration: 280,
        useNativeDriver: true,
      }),
    ]).start(() => {
      resetAll();
      onClose();
    });
  }

  function resetAll() {
    stopInterval();
    setTimerState('idle');
    timerStateRef.current = 'idle';
    setElapsed(0);
    setSelectedBean(null);
    setDoseIn('');
    setYieldOut('');
    setGrindSize('');
    setTargetTime(27);
    setPreInfusion(false);
    setSelectedGearIds([]);
    setTastingNote('');
    setManualMode(false);
    setManualSeconds('');
    setBeanSheetOpen(false);
    setGearSheetOpen(false);
  }

  // ── Sub-sheets ──
  function openBeanSheet() {
    setBeanSheetOpen(true);
    beanSheetAnim.setValue(SUB_SHEET_MAX);
    Animated.timing(beanSheetAnim, {
      toValue: 0,
      duration: 300,
      useNativeDriver: true,
    }).start();
  }

  function closeBeanSheet() {
    Animated.timing(beanSheetAnim, {
      toValue: SUB_SHEET_MAX,
      duration: 250,
      useNativeDriver: true,
    }).start(() => setBeanSheetOpen(false));
  }

  function openGearSheet() {
    setGearSheetOpen(true);
    gearSheetAnim.setValue(SUB_SHEET_MAX);
    Animated.timing(gearSheetAnim, {
      toValue: 0,
      duration: 300,
      useNativeDriver: true,
    }).start();
  }

  function closeGearSheet() {
    Animated.timing(gearSheetAnim, {
      toValue: SUB_SHEET_MAX,
      duration: 250,
      useNativeDriver: true,
    }).start(() => setGearSheetOpen(false));
  }

  // ── Timer logic ──
  function stopInterval() {
    if (intervalRef.current) {
      clearInterval(intervalRef.current);
      intervalRef.current = null;
    }
  }

  function startExtraction() {
    runStartRef.current = Date.now();
    timerStateRef.current = 'running';
    setTimerState('running');

    intervalRef.current = setInterval(() => {
      if (timerStateRef.current === 'running') {
        setElapsed((Date.now() - runStartRef.current) / 1000);
      }
    }, 50);
  }

  function stopExtraction() {
    stopInterval();
    timerStateRef.current = 'done';
    setTimerState('done');
  }

  function resetExtraction() {
    stopInterval();
    setElapsed(0);
    timerStateRef.current = 'idle';
    setTimerState('idle');
    setTastingNote('');
  }

  // ── Log manual time ──
  function logManual() {
    const secs = parseFloat(manualSeconds);
    if (!isNaN(secs) && secs > 0) {
      setElapsed(secs);
      timerStateRef.current = 'done';
      setTimerState('done');
    }
  }

  // ── Save ──
  async function saveExtraction() {
    if (!selectedBean) {
      Alert.alert('Select beans', 'Please select the beans used for this extraction.');
      return;
    }
    const dose = parseFloat(doseIn);
    const yld = parseFloat(yieldOut);
    const grind = parseFloat(grindSize);
    if (!dose || dose <= 0) {
      Alert.alert('Dose required', 'Enter the dose in grams.');
      return;
    }
    if (!yld || yld <= 0) {
      Alert.alert('Yield required', 'Enter the yield out in grams.');
      return;
    }
    if (!grind || grind <= 0) {
      Alert.alert('Grind size required', 'Enter the grind size.');
      return;
    }
    if (elapsed <= 0) {
      Alert.alert('No time recorded', 'Run the timer or enter a time manually.');
      return;
    }

    setSaving(true);
    try {
      const result = await extractionApi.create({
        bean_id: selectedBean.id,
        dose_in: dose,
        yield_out: yld,
        time: Math.round(elapsed * 10) / 10,
        target_time: targetTime,
        grind_size: grind,
        gear_ids: selectedGearIds,
        pre_infusion: preInfusion,
        tasting_note: tastingNote.trim() || null,
      });
      addExtraction(result);
      dismiss();
    } catch {
      Alert.alert('Save failed', 'Could not save the extraction. Please try again.');
    } finally {
      setSaving(false);
    }
  }

  // ── Computed ──
  const doseNum = parseFloat(doseIn);
  const yieldNum = parseFloat(yieldOut);
  const ratio =
    doseNum > 0 && yieldNum > 0
      ? `1:${(yieldNum / doseNum).toFixed(1)}`
      : null;

  const selectedGear = gear.filter(g => selectedGearIds.includes(g.id));
  const gearLabel = selectedGear.length
    ? selectedGear.map(g => g.name).join(' · ')
    : 'Add gear';

  const isIdle = timerState === 'idle';
  const isRunning = timerState === 'running';
  const isDone = timerState === 'done';
  const fieldsLocked = isRunning;

  const subSheetActive = beanSheetOpen || gearSheetOpen;

  // ── Target time stepper ──
  function incrTarget() {
    setTargetTime(t => Math.min(t + 1, 90));
  }
  function decrTarget() {
    setTargetTime(t => Math.max(t - 1, 10));
  }

  // ── Cleanup ──
  useEffect(() => {
    return () => stopInterval();
  }, []);

  return (
    <Modal
      visible={visible}
      transparent
      animationType="none"
      statusBarTranslucent
      onRequestClose={dismiss}
    >
      <View style={styles.modalRoot}>
        {/* Backdrop: outside KAV so it covers full screen including keyboard area */}
        <Animated.View style={[styles.backdrop, { opacity: backdropAnim }]}>
          <Pressable style={StyleSheet.absoluteFill} onPress={dismiss} />
        </Animated.View>

        <Animated.View
          style={[styles.sheet, { transform: [{ translateY: sheetAnim }] }]}
        >
          {/* Handle */}
          <View style={styles.handle} />

          {/* Header */}
          <View style={styles.modalHeader}>
            <Text style={styles.modalTitle}>New extraction</Text>
            {!isRunning && (
              <Pressable style={styles.closeBtn} onPress={dismiss}>
                <Feather name="x" size={16} color={palette.espresso500} />
              </Pressable>
            )}
          </View>

          <KeyboardAwareScrollView
            style={styles.scroll}
            contentContainerStyle={styles.scrollContent}
            keyboardShouldPersistTaps="handled"
            bottomOffset={24}
          >
            {/* Bean selector */}
            <Pressable
              style={[
                styles.beanSelector,
                selectedBean
                  ? styles.beanSelectorFilled
                  : styles.beanSelectorEmpty,
              ]}
              onPress={!fieldsLocked ? openBeanSheet : undefined}
            >
              <View
                style={[
                  styles.beanSelectorIcon,
                  {
                    backgroundColor: selectedBean
                      ? palette.espresso700
                      : palette.cream300,
                  },
                ]}
              >
                <View
                  style={[
                    styles.beanTearMd,
                    {
                      backgroundColor: selectedBean
                        ? roastColor(selectedBean.roast_level)
                        : palette.cream500,
                    },
                  ]}
                />
              </View>
              <View style={styles.beanSelectorInfo}>
                {selectedBean ? (
                  <>
                    <Text
                      style={[
                        styles.beanSelectorName,
                        { color: palette.cream100 },
                      ]}
                    >
                      {selectedBean.name}
                    </Text>
                    {(selectedBean.roaster || selectedBean.roast_level) && (
                      <Text style={[styles.beanSelectorMeta, { color: palette.cream500 }]}>
                        {[selectedBean.roaster, selectedBean.roast_level ? roastLabel(selectedBean.roast_level) : null]
                          .filter(Boolean)
                          .join(' · ')}
                      </Text>
                    )}
                  </>
                ) : (
                  <Text style={styles.beanSelectorPlaceholder}>Select beans</Text>
                )}
              </View>
              <Feather
                name="chevron-right"
                size={16}
                color={selectedBean ? palette.cream500 : palette.espresso400}
              />
            </Pressable>

            {/* Timer ring */}
            <View style={styles.ringContainer}>
              <TimerRing
                timerState={timerState}
                elapsed={elapsed}
                targetTime={targetTime}
              />
              <View style={styles.ringCenter}>
                {isIdle ? (
                  <>
                    <Text style={styles.ringReadyLabel}>ready</Text>
                    <Text style={styles.ringIdleTime}>--:--</Text>
                  </>
                ) : (
                  <>
                    <Text style={[styles.ringTime, { color: isDone ? zoneColor(elapsed, targetTime) : palette.espresso800 }]}>
                      {fmtTime(elapsed)}
                    </Text>
                    <Text style={[styles.ringZoneLabel, { color: zoneColor(elapsed, targetTime) }]}>
                      {zoneLabel(elapsed, targetTime).toUpperCase()}
                    </Text>
                  </>
                )}
              </View>
              {isDone && (
                <View style={[styles.zoneBadge, { backgroundColor: zoneColor(elapsed, targetTime) }]}>
                  <Text style={styles.zoneBadgeText}>
                    {zoneLabel(elapsed, targetTime)} · {Math.round(elapsed)}s
                  </Text>
                </View>
              )}
            </View>

            {/* Aim for + pre-infusion row (idle only) */}
            {isIdle && (
              <>
                <View style={styles.aimRow}>
                  <View style={styles.aimStepper}>
                    <Text style={styles.aimLabel}>Aim for</Text>
                    <View style={styles.stepperRow}>
                      <Pressable style={styles.stepperBtn} onPress={decrTarget}>
                        <Text style={styles.stepperBtnText}>−</Text>
                      </Pressable>
                      <Text style={styles.stepperVal}>{targetTime}s</Text>
                      <Pressable style={styles.stepperBtn} onPress={incrTarget}>
                        <Text style={styles.stepperBtnText}>+</Text>
                      </Pressable>
                    </View>
                  </View>
                  <Pressable
                    style={[
                      styles.preInfToggle,
                      preInfusion ? styles.preInfToggleOn : styles.preInfToggleOff,
                    ]}
                    onPress={() => setPreInfusion(p => !p)}
                  >
                    <View
                      style={[
                        styles.preInfDot,
                        { backgroundColor: preInfusion ? palette.caramel400 : palette.cream500 },
                      ]}
                    />
                    <Text
                      style={[
                        styles.preInfText,
                        { color: preInfusion ? palette.caramel500 : palette.espresso400 },
                      ]}
                    >
                      Pre-infusion
                    </Text>
                  </Pressable>
                </View>
                <View style={styles.divider} />
              </>
            )}

            {/* Stats row */}
            <View style={styles.statsRow}>
              <StatField
                label="DOSE IN"
                unit="g"
                value={doseIn}
                onChangeText={setDoseIn}
                editable={!fieldsLocked}
              />
              <View style={styles.statDivider} />
              <StatField
                label="YIELD OUT"
                unit="g"
                value={yieldOut}
                onChangeText={setYieldOut}
                editable={!fieldsLocked}
              />
              <View style={styles.statDivider} />
              <StatField
                label="GRIND"
                unit=""
                value={grindSize}
                onChangeText={setGrindSize}
                editable={!fieldsLocked}
              />
            </View>
            {ratio && (
              <Text style={styles.ratioText}>{ratio} ratio</Text>
            )}
            <View style={styles.divider} />

            {/* Gear row */}
            <Pressable
              style={styles.gearRow}
              onPress={!isRunning ? openGearSheet : undefined}
            >
              <View>
                <Text style={styles.gearLabel}>GEAR USED</Text>
                <Text
                  style={[
                    styles.gearValue,
                    { color: selectedGear.length ? palette.espresso800 : palette.espresso400 },
                  ]}
                  numberOfLines={2}
                >
                  {gearLabel}
                </Text>
              </View>
              {!isRunning && (
                <Feather name="chevron-right" size={16} color={palette.espresso400} />
              )}
            </Pressable>
            <View style={styles.divider} />

            {/* Tasting note (done state) */}
            {isDone && (
              <View style={styles.noteContainer}>
                <Text style={styles.noteLabel}>TASTING NOTE</Text>
                <TextInput
                  style={styles.noteInput}
                  value={tastingNote}
                  onChangeText={setTastingNote}
                  multiline
                  placeholder="How did it taste? Sour, sweet, balanced…"
                  placeholderTextColor={palette.cream500}
                />
              </View>
            )}

            {/* CTA buttons — inline so no layout jump between modes */}
            <View style={[styles.ctaSection, { paddingBottom: insets.bottom + 12 }]}>
              {isIdle && !manualMode && (
                <>
                  <Pressable style={styles.startBtn} onPress={startExtraction}>
                    <Feather name="play" size={18} color="#fff" />
                    <Text style={styles.ctaBtnText}>Start extraction</Text>
                  </Pressable>
                  <Pressable onPress={() => setManualMode(true)}>
                    <Text style={styles.manualToggle}>Add time manually</Text>
                  </Pressable>
                </>
              )}
              {isIdle && manualMode && (
                <>
                  <View style={styles.manualRow}>
                    <TextInput
                      style={styles.manualInput}
                      value={manualSeconds}
                      onChangeText={setManualSeconds}
                      keyboardType="decimal-pad"
                      placeholder="seconds"
                      placeholderTextColor={palette.cream500}
                    />
                    <Pressable style={styles.manualLogBtn} onPress={logManual}>
                      <Text style={styles.manualLogText}>Log</Text>
                    </Pressable>
                  </View>
                  <Pressable onPress={() => setManualMode(false)}>
                    <Text style={styles.manualToggle}>← Use live timer instead</Text>
                  </Pressable>
                </>
              )}
              {isRunning && (
                <Pressable style={styles.stopBtn} onPress={stopExtraction}>
                  <Feather name="square" size={18} color="#fff" />
                  <Text style={styles.ctaBtnText}>Stop</Text>
                </Pressable>
              )}
              {isDone && (
                <View style={styles.doneRow}>
                  <Pressable style={styles.resetBtn} onPress={resetExtraction}>
                    <Feather name="refresh-cw" size={20} color={palette.espresso500} />
                  </Pressable>
                  <Pressable
                    style={[styles.saveBtn, saving && styles.saveBtnDisabled]}
                    onPress={saveExtraction}
                    disabled={saving}
                  >
                    <Text style={styles.ctaBtnText}>
                      {saving ? 'Saving…' : 'Log extraction'}
                    </Text>
                  </Pressable>
                </View>
              )}
            </View>
          </KeyboardAwareScrollView>
        </Animated.View>

        {/* Sub-sheet backdrop */}
        {subSheetActive && (
          <Pressable
            style={[StyleSheet.absoluteFill, styles.subBackdrop]}
            onPress={beanSheetOpen ? closeBeanSheet : closeGearSheet}
          />
        )}

        {/* Bean sub-sheet */}
        {beanSheetOpen && (
          <BeanSheet
            beans={beans}
            loading={beansLoading}
            selectedId={selectedBean?.id ?? null}
            onSelect={setSelectedBean}
            onClose={closeBeanSheet}
            translateY={beanSheetAnim}
            onRefresh={refreshBeans}
          />
        )}

        {/* Gear sub-sheet */}
        {gearSheetOpen && (
          <GearSheet
            gear={gear}
            stations={stations}
            loading={gearLoading}
            selectedIds={selectedGearIds}
            onSave={setSelectedGearIds}
            onClose={closeGearSheet}
            translateY={gearSheetAnim}
            onRefresh={refreshGear}
          />
        )}
      </View>
    </Modal>
  );
}

// ─── StatField ────────────────────────────────────────────────────────────────

function StatField({
  label,
  unit,
  value,
  onChangeText,
  editable,
}: {
  label: string;
  unit: string;
  value: string;
  onChangeText: (v: string) => void;
  editable: boolean;
}) {
  const [focused, setFocused] = useState(false);
  const isEmpty = !value;

  return (
    <View style={styles.statField}>
      <Text style={styles.statLabel}>{label}</Text>
      <View style={styles.statValueRow}>
        {editable ? (
          <TextInput
            style={[
              styles.statInput,
              focused && styles.statInputFocused,
              isEmpty && styles.statInputEmpty,
            ]}
            value={value}
            onChangeText={onChangeText}
            keyboardType="decimal-pad"
            placeholder="—"
            placeholderTextColor={palette.cream500}
            onFocus={() => setFocused(true)}
            onBlur={() => setFocused(false)}
          />
        ) : (
          <Text style={[styles.statInputLocked, isEmpty && { color: palette.cream500 }]}>
            {isEmpty ? '—' : value}
          </Text>
        )}
        {!isEmpty && !!unit && <Text style={styles.statUnit}>{unit}</Text>}
      </View>
    </View>
  );
}

// ─── Styles ───────────────────────────────────────────────────────────────────

const styles = StyleSheet.create({
  modalRoot: {
    flex: 1,
  },
  backdrop: {
    ...StyleSheet.absoluteFillObject,
    backgroundColor: 'rgba(28,15,7,0.50)',
  },
  sheet: {
    position: 'absolute',
    bottom: 0,
    left: 0,
    right: 0,
    height: SHEET_HEIGHT,
    backgroundColor: palette.cream100,
    borderTopLeftRadius: 32,
    borderTopRightRadius: 32,
    shadowColor: '#1C0F07',
    shadowOffset: { width: 0, height: -8 },
    shadowOpacity: 0.18,
    shadowRadius: 24,
    elevation: 16,
  },
  handle: {
    width: 38,
    height: 4,
    backgroundColor: palette.cream500,
    borderRadius: radii.full,
    alignSelf: 'center',
    marginTop: 14,
  },
  modalHeader: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    paddingHorizontal: 22,
    paddingTop: 16,
    paddingBottom: 4,
  },
  modalTitle: {
    fontFamily: 'PlayfairDisplay_700Bold_Italic',
    fontSize: 22,
    color: palette.espresso800,
  },
  closeBtn: {
    width: 32,
    height: 32,
    borderRadius: radii.full,
    backgroundColor: palette.cream300,
    alignItems: 'center',
    justifyContent: 'center',
  },
  scroll: {
    flex: 1,
  },
  scrollContent: {
    paddingTop: 0,
  },

  // Bean selector
  beanSelector: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: 10,
    marginHorizontal: 20,
    marginTop: 16,
    padding: 12,
    paddingHorizontal: 16,
    borderRadius: 16,
    borderWidth: 1.5,
  },
  beanSelectorFilled: {
    backgroundColor: palette.espresso800,
    borderColor: palette.espresso800,
  },
  beanSelectorEmpty: {
    backgroundColor: palette.cream200,
    borderColor: palette.cream400,
  },
  beanSelectorIcon: {
    width: 28,
    height: 28,
    borderRadius: 8,
    alignItems: 'center',
    justifyContent: 'center',
  },
  beanTearMd: {
    width: 14,
    height: 18,
    borderTopLeftRadius: 7,
    borderTopRightRadius: 7,
    borderBottomLeftRadius: 5,
    borderBottomRightRadius: 5,
  },
  beanSelectorInfo: {
    flex: 1,
  },
  beanSelectorName: {
    fontFamily: 'DMSans_700Bold',
    fontSize: 14,
  },
  beanSelectorMeta: {
    fontFamily: 'DMSans_400Regular',
    fontSize: 12,
    marginTop: 1,
  },
  beanSelectorPlaceholder: {
    fontFamily: 'DMSans_400Regular',
    fontSize: 14,
    color: palette.espresso400,
  },

  // Ring
  ringContainer: {
    alignItems: 'center',
    paddingTop: 16,
    paddingBottom: 10,
  },
  ringCenter: {
    position: 'absolute',
    top: 16,
    left: 0,
    right: 0,
    height: RING_SIZE,
    alignItems: 'center',
    justifyContent: 'center',
  },
  ringReadyLabel: {
    fontFamily: 'JetBrainsMono_500Medium',
    fontSize: 14,
    letterSpacing: 1,
    color: palette.cream500,
    marginBottom: 4,
  },
  ringIdleTime: {
    fontFamily: 'JetBrainsMono_600SemiBold',
    fontSize: 44,
    letterSpacing: -2,
    color: palette.cream400,
  },
  ringTime: {
    fontFamily: 'JetBrainsMono_600SemiBold',
    fontSize: 52,
    letterSpacing: -2,
  },
  ringZoneLabel: {
    fontFamily: 'DMSans_700Bold',
    fontSize: 12,
    letterSpacing: 0.05 * 12,
    marginTop: 4,
  },
  zoneBadge: {
    borderRadius: radii.full,
    paddingVertical: 6,
    paddingHorizontal: 20,
    marginTop: 12,
  },
  zoneBadgeText: {
    fontFamily: 'DMSans_700Bold',
    fontSize: 13,
    color: '#fff',
    letterSpacing: 0.04 * 13,
  },

  // Aim row
  aimRow: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    paddingHorizontal: 24,
    paddingBottom: 24,
  },
  aimStepper: {
    alignItems: 'flex-start',
    gap: 6,
  },
  aimLabel: {
    fontFamily: 'DMSans_500Medium',
    fontSize: 13,
    color: palette.espresso500,
  },
  stepperRow: {
    flexDirection: 'row',
    alignItems: 'center',
    backgroundColor: palette.cream300,
    borderRadius: radii.full,
    overflow: 'hidden',
  },
  stepperBtn: {
    width: 34,
    height: 32,
    alignItems: 'center',
    justifyContent: 'center',
  },
  stepperBtnText: {
    fontFamily: 'DMSans_400Regular',
    fontSize: 18,
    color: palette.espresso500,
    lineHeight: 20,
  },
  stepperVal: {
    fontFamily: 'JetBrainsMono_600SemiBold',
    fontSize: 14,
    color: palette.espresso800,
    minWidth: 32,
    textAlign: 'center',
  },
  preInfToggle: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: 6,
    borderWidth: 1.5,
    borderRadius: radii.full,
    paddingVertical: 6,
    paddingHorizontal: 12,
  },
  preInfToggleOn: {
    backgroundColor: palette.caramel100,
    borderColor: palette.caramel400,
  },
  preInfToggleOff: {
    backgroundColor: 'transparent',
    borderColor: palette.cream400,
  },
  preInfDot: {
    width: 8,
    height: 8,
    borderRadius: radii.full,
  },
  preInfText: {
    fontFamily: 'DMSans_700Bold',
    fontSize: 12,
  },
  divider: {
    height: 1,
    backgroundColor: palette.cream300,
    marginHorizontal: 22,
  },

  // Stats
  statsRow: {
    flexDirection: 'row',
    paddingHorizontal: 24,
    paddingVertical: 24,
  },
  statField: {
    flex: 1,
    alignItems: 'center',
  },
  statLabel: {
    fontFamily: 'DMSans_700Bold',
    fontSize: 11,
    letterSpacing: 0.07 * 11,
    color: palette.espresso400,
    marginBottom: 6,
    textTransform: 'uppercase',
  },
  statValueRow: {
    flexDirection: 'row',
    alignItems: 'flex-end',
    gap: 2,
  },
  statInput: {
    fontFamily: 'JetBrainsMono_600SemiBold',
    fontSize: 28,
    letterSpacing: -1,
    color: palette.espresso800,
    padding: 0,
    minWidth: 40,
    textAlign: 'center',
    borderBottomWidth: 1.5,
    borderBottomColor: 'transparent',
  },
  statInputFocused: {
    borderBottomColor: palette.caramel400,
  },
  statInputEmpty: {
    color: palette.cream500,
  },
  statInputLocked: {
    fontFamily: 'JetBrainsMono_600SemiBold',
    fontSize: 28,
    letterSpacing: -1,
    color: palette.espresso800,
  },
  statUnit: {
    fontFamily: 'JetBrainsMono_400Regular',
    fontSize: 14,
    color: palette.espresso400,
    marginBottom: 4,
  },
  statDivider: {
    width: 1,
    alignSelf: 'stretch',
    backgroundColor: palette.cream400,
    marginVertical: 4,
  },
  ratioText: {
    fontFamily: 'JetBrainsMono_600SemiBold',
    fontSize: 13,
    color: palette.caramel500,
    textAlign: 'center',
    marginTop: -12,
    marginBottom: 16,
  },

  // Gear
  gearRow: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    paddingHorizontal: 24,
    paddingVertical: 18,
  },
  gearLabel: {
    fontFamily: 'DMSans_700Bold',
    fontSize: 12,
    letterSpacing: 0.07 * 12,
    color: palette.espresso400,
    textTransform: 'uppercase',
    marginBottom: 4,
  },
  gearValue: {
    fontFamily: 'DMSans_700Bold',
    fontSize: 14,
  },

  // Tasting note
  noteContainer: {
    paddingHorizontal: 22,
    paddingVertical: 18,
  },
  noteLabel: {
    fontFamily: 'DMSans_700Bold',
    fontSize: 12,
    letterSpacing: 0.07 * 12,
    color: palette.espresso400,
    textTransform: 'uppercase',
    marginBottom: 10,
  },
  noteInput: {
    minHeight: 72,
    padding: 12,
    paddingHorizontal: 14,
    backgroundColor: palette.cream200,
    borderWidth: 1.5,
    borderColor: palette.cream400,
    borderRadius: 14,
    fontFamily: 'DMSans_400Regular',
    fontSize: 14,
    lineHeight: 21,
    color: palette.espresso800,
    textAlignVertical: 'top',
  },

  // CTA
  ctaSection: {
    paddingHorizontal: 22,
    paddingTop: 12,
    gap: 10,
    borderTopWidth: 1,
    borderTopColor: palette.cream300,
  },
  startBtn: {
    height: 56,
    borderRadius: radii.full,
    backgroundColor: palette.matcha500,
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'center',
    gap: 8,
    shadowColor: palette.matcha500,
    shadowOffset: { width: 0, height: 6 },
    shadowOpacity: 0.38,
    shadowRadius: 12,
    elevation: 4,
  },
  stopBtn: {
    height: 56,
    borderRadius: radii.full,
    backgroundColor: palette.error500,
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'center',
    gap: 8,
    shadowColor: palette.error500,
    shadowOffset: { width: 0, height: 6 },
    shadowOpacity: 0.35,
    shadowRadius: 12,
    elevation: 4,
  },
  ctaBtnText: {
    fontFamily: 'DMSans_700Bold',
    fontSize: 16,
    color: '#fff',
  },
  doneRow: {
    flexDirection: 'row',
    gap: 10,
  },
  resetBtn: {
    width: 56,
    height: 56,
    borderRadius: radii.full,
    backgroundColor: palette.cream300,
    alignItems: 'center',
    justifyContent: 'center',
  },
  saveBtn: {
    flex: 1,
    height: 56,
    borderRadius: radii.full,
    backgroundColor: palette.caramel400,
    alignItems: 'center',
    justifyContent: 'center',
    shadowColor: palette.caramel400,
    shadowOffset: { width: 0, height: 6 },
    shadowOpacity: 0.38,
    shadowRadius: 10,
    elevation: 4,
  },
  saveBtnDisabled: {
    opacity: 0.6,
  },
  manualRow: {
    flexDirection: 'row',
    gap: 10,
    alignItems: 'center',
  },
  manualInput: {
    flex: 1,
    height: 46,
    borderRadius: 12,
    borderWidth: 1.5,
    borderColor: palette.cream400,
    backgroundColor: palette.cream200,
    paddingHorizontal: 14,
    fontFamily: 'JetBrainsMono_600SemiBold',
    fontSize: 16,
    color: palette.espresso800,
    textAlign: 'center',
  },
  manualLogBtn: {
    height: 46,
    paddingHorizontal: 20,
    borderRadius: radii.full,
    backgroundColor: palette.espresso800,
    alignItems: 'center',
    justifyContent: 'center',
  },
  manualLogText: {
    fontFamily: 'DMSans_700Bold',
    fontSize: 14,
    color: '#fff',
  },
  manualToggle: {
    fontFamily: 'DMSans_500Medium',
    fontSize: 13,
    color: palette.espresso400,
    textAlign: 'center',
  },

  // Sub-sheet backdrop
  subBackdrop: {
    backgroundColor: 'rgba(28,15,7,0.30)',
  },

  // Sub-sheet shared
  subSheet: {
    position: 'absolute',
    bottom: 0,
    left: 0,
    right: 0,
    maxHeight: SUB_SHEET_MAX,
    backgroundColor: palette.cream100,
    borderTopLeftRadius: 28,
    borderTopRightRadius: 28,
    shadowColor: '#1C0F07',
    shadowOffset: { width: 0, height: -4 },
    shadowOpacity: 0.14,
    shadowRadius: 20,
    elevation: 20,
  },
  subHandle: {
    width: 38,
    height: 4,
    backgroundColor: palette.cream500,
    borderRadius: radii.full,
    alignSelf: 'center',
    marginTop: 12,
  },
  subHeader: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    paddingHorizontal: 20,
    paddingVertical: 14,
  },
  subTitle: {
    fontFamily: 'PlayfairDisplay_700Bold',
    fontSize: 18,
    color: palette.espresso800,
  },
  subClose: {
    width: 30,
    height: 30,
    borderRadius: radii.full,
    backgroundColor: palette.cream300,
    alignItems: 'center',
    justifyContent: 'center',
  },
  subList: {
    flex: 1,
  },
  subListContent: {
    paddingHorizontal: 16,
    paddingBottom: 32,
    gap: 8,
  },

  // Bean sheet rows
  beanRow: {
    flexDirection: 'row',
    alignItems: 'center',
    gap: 12,
    padding: 13,
    paddingHorizontal: 15,
    borderRadius: 16,
    backgroundColor: palette.cream200,
    borderWidth: 1.5,
    borderColor: palette.cream400,
  },
  beanRowSel: {
    backgroundColor: palette.espresso800,
    borderColor: palette.espresso800,
  },
  beanIconBox: {
    width: 36,
    height: 36,
    borderRadius: 8,
    alignItems: 'center',
    justifyContent: 'center',
  },
  beanTear: {
    width: 16,
    height: 20,
    borderTopLeftRadius: 8,
    borderTopRightRadius: 8,
    borderBottomLeftRadius: 6,
    borderBottomRightRadius: 6,
  },
  beanInfo: {
    flex: 1,
    gap: 2,
  },
  emptyHint: {
    fontFamily: 'DMSans_400Regular',
    fontSize: 14,
    color: palette.espresso400,
    textAlign: 'center',
    paddingVertical: 24,
  },
  beanRowName: {
    fontFamily: 'DMSans_700Bold',
    fontSize: 14,
  },
  beanRowMeta: {
    fontFamily: 'DMSans_400Regular',
    fontSize: 12,
  },
  addBeanBtn: {
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'center',
    gap: 6,
    padding: 14,
    borderRadius: 16,
    borderWidth: 1.5,
    borderStyle: 'dashed',
    borderColor: palette.cream400,
  },
  addBeanText: {
    fontFamily: 'DMSans_500Medium',
    fontSize: 13,
    color: palette.espresso400,
  },

  // Gear sheet
  gearSectionLabel: {
    fontFamily: 'DMSans_700Bold',
    fontSize: 11,
    letterSpacing: 0.07 * 11,
    textTransform: 'uppercase',
    color: palette.espresso400,
    marginTop: 4,
    marginBottom: 4,
  },
  stationRow: {
    flexDirection: 'row',
    gap: 8,
    flexWrap: 'wrap',
    marginBottom: 8,
  },
  stationBtn: {
    paddingVertical: 8,
    paddingHorizontal: 14,
    borderRadius: radii.full,
    backgroundColor: palette.cream300,
    borderWidth: 1,
    borderColor: palette.cream400,
  },
  stationBtnText: {
    fontFamily: 'DMSans_700Bold',
    fontSize: 13,
    color: palette.espresso800,
  },
  gearItemRow: {
    flexDirection: 'row',
    alignItems: 'center',
    padding: 13,
    paddingHorizontal: 15,
    borderRadius: 16,
    backgroundColor: palette.cream200,
    borderWidth: 1.5,
    borderColor: palette.cream400,
  },
  gearItemRowSel: {
    backgroundColor: palette.espresso800,
    borderColor: palette.espresso800,
  },
  gearInfo: {
    flex: 1,
    gap: 2,
  },
  gearItemName: {
    fontFamily: 'DMSans_700Bold',
    fontSize: 14,
  },
  gearItemType: {
    fontFamily: 'DMSans_400Regular',
    fontSize: 12,
    textTransform: 'capitalize',
  },
  gearFooter: {
    padding: 16,
    paddingHorizontal: 16,
    borderTopWidth: 1,
    borderTopColor: palette.cream300,
  },
  gearSaveBtn: {
    height: 50,
    borderRadius: radii.full,
    backgroundColor: palette.espresso800,
    alignItems: 'center',
    justifyContent: 'center',
  },
  gearSaveBtnText: {
    fontFamily: 'DMSans_700Bold',
    fontSize: 15,
    color: palette.cream100,
  },
});
