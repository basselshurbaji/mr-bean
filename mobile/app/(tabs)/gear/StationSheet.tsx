import {
  ActivityIndicator,
  Animated,
  KeyboardAvoidingView,
  Modal,
  Platform,
  Pressable,
  ScrollView,
  StyleSheet,
  Text,
  TextInput,
  View,
} from 'react-native';
import { useEffect, useRef, useState } from 'react';
import { colors, palette, radii, spacing } from '@/src/theme';
import { gearApi, GearItem, Station } from '@/src/api/gear';
import GearIcon from '@/src/components/GearIcon';

interface Props {
  gear: GearItem[];
  editStation?: Station;
  onClose: () => void;
  onSaved: (station: Station) => void;
  onDeleted?: (id: string) => void;
}

export default function StationSheet({ gear, editStation, onClose, onSaved, onDeleted }: Props) {
  const translateY = useRef(new Animated.Value(600)).current;
  const backdropOpacity = useRef(new Animated.Value(0)).current;

  const isEdit = !!editStation;

  const initialSelectedIds = editStation?.gear.map(g => g.id) ?? [];

  const [name, setName] = useState(editStation?.name ?? '');
  const [selectedIds, setSelectedIds] = useState<string[]>(initialSelectedIds);
  const [focused, setFocused] = useState(false);
  const [saving, setSaving] = useState(false);
  const [deleting, setDeleting] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const sortedGear = [...gear].sort((a, b) => {
    const aIdx = initialSelectedIds.indexOf(a.id);
    const bIdx = initialSelectedIds.indexOf(b.id);
    const aSelected = aIdx !== -1;
    const bSelected = bIdx !== -1;
    if (aSelected && bSelected) return aIdx - bIdx;
    if (aSelected) return -1;
    if (bSelected) return 1;
    return 0;
  });

  useEffect(() => {
    Animated.parallel([
      Animated.timing(translateY, { toValue: 0, duration: 300, useNativeDriver: true }),
      Animated.timing(backdropOpacity, { toValue: 1, duration: 300, useNativeDriver: true }),
    ]).start();
  }, []);

  function dismiss() {
    Animated.parallel([
      Animated.timing(translateY, { toValue: 600, duration: 260, useNativeDriver: true }),
      Animated.timing(backdropOpacity, { toValue: 0, duration: 260, useNativeDriver: true }),
    ]).start(() => onClose());
  }

  function toggleGear(id: string) {
    setSelectedIds(prev =>
      prev.includes(id) ? prev.filter(x => x !== id) : [...prev, id],
    );
  }

  async function save() {
    if (saving || !name.trim()) return;
    setSaving(true);
    setError(null);
    try {
      const body = { name: name.trim(), gear_ids: selectedIds };
      const result = isEdit
        ? await gearApi.updateStation(editStation!.id, body)
        : await gearApi.createStation(body);
      onSaved(result);
    } catch (e) {
      setError(e instanceof Error ? e.message : 'Something went wrong.');
      setSaving(false);
    }
  }

  async function deleteStation() {
    if (deleting || !editStation) return;
    setDeleting(true);
    try {
      await gearApi.deleteStation(editStation.id);
      onDeleted?.(editStation.id);
    } catch {
      setError('Failed to delete station.');
      setDeleting(false);
    }
  }

  const canSave = name.trim().length > 0;

  return (
    <Modal transparent animationType="none" onRequestClose={dismiss}>
      <View style={styles.overlay}>
        <Animated.View style={[styles.backdrop, { opacity: backdropOpacity }]}>
          <Pressable style={StyleSheet.absoluteFill} onPress={dismiss} />
        </Animated.View>

        <Animated.View style={[styles.sheet, { transform: [{ translateY }] }]}>
          <View style={styles.handle} />

          <KeyboardAvoidingView
            behavior={Platform.OS === 'ios' ? 'padding' : 'height'}
          >
            <ScrollView
              contentContainerStyle={styles.content}
              showsVerticalScrollIndicator={false}
              keyboardShouldPersistTaps="handled"
            >
              <Text style={styles.title}>
                {isEdit ? 'Edit station' : 'New station'}
              </Text>

              {/* Station name */}
              <View style={styles.fieldWrap}>
                <Text style={styles.fieldLabel}>
                  Station name <Text style={styles.required}>*</Text>
                </Text>
                <TextInput
                  style={[styles.input, focused && styles.inputFocused]}
                  placeholder="e.g. Morning routine"
                  placeholderTextColor={colors.fgDisabled}
                  value={name}
                  onChangeText={setName}
                  onFocus={() => setFocused(true)}
                  onBlur={() => setFocused(false)}
                  returnKeyType="done"
                />
              </View>

              <Text style={styles.hint}>
                Stations pre-select tools when logging — you can always tweak before saving a shot.
              </Text>

              {/* Gear list */}
              {sortedGear.length > 0 && (
                <View style={styles.gearList}>
                  {sortedGear.map(item => {
                    const selected = selectedIds.includes(item.id);
                    const sub = [item.brand, item.model].filter(Boolean).join(' · ');
                    return (
                      <Pressable
                        key={item.id}
                        style={[styles.gearRow, selected && styles.gearRowSelected]}
                        onPress={() => toggleGear(item.id)}
                      >
                        <View style={[styles.gearBubble, selected && styles.gearBubbleSelected]}>
                          <GearIcon
                            typeId={item.type_id}
                            size={22}
                            color={selected ? palette.cream100 : palette.espresso800}
                          />
                        </View>
                        <View style={styles.gearText}>
                          <Text style={[styles.gearName, selected && styles.gearNameSelected]}>
                            {item.name}
                          </Text>
                          {sub ? (
                            <Text style={[styles.gearSub, selected && styles.gearSubSelected]}>
                              {sub}
                            </Text>
                          ) : null}
                        </View>
                        <View style={[styles.toggle, selected && styles.toggleSelected]}>
                          {selected && (
                            <Text style={styles.toggleCheck}>✓</Text>
                          )}
                        </View>
                      </Pressable>
                    );
                  })}
                </View>
              )}

              {error && (
                <View style={styles.errorBanner}>
                  <Text style={styles.errorText}>{error}</Text>
                </View>
              )}

              <Pressable
                style={({ pressed }) => [
                  styles.cta,
                  (!canSave || saving) && styles.ctaDisabled,
                  pressed && canSave && !saving && styles.ctaPressed,
                ]}
                onPress={save}
                disabled={!canSave || saving}
              >
                {saving ? (
                  <ActivityIndicator color={palette.cream100} />
                ) : (
                  <Text style={styles.ctaLabel}>
                    {isEdit
                      ? 'Save changes'
                      : `Create station · ${selectedIds.length} item${selectedIds.length !== 1 ? 's' : ''}`}
                  </Text>
                )}
              </Pressable>

              {isEdit && (
                <Pressable
                  style={({ pressed }) => [
                    styles.deleteBtn,
                    pressed && { opacity: 0.75 },
                    deleting && { opacity: 0.38 },
                  ]}
                  onPress={deleteStation}
                  disabled={deleting}
                >
                  {deleting ? (
                    <ActivityIndicator color={palette.error500} />
                  ) : (
                    <Text style={styles.deleteBtnLabel}>Delete station</Text>
                  )}
                </Pressable>
              )}
            </ScrollView>
          </KeyboardAvoidingView>
        </Animated.View>
      </View>
    </Modal>
  );
}

const styles = StyleSheet.create({
  overlay: { flex: 1, justifyContent: 'flex-end' },
  backdrop: {
    ...StyleSheet.absoluteFillObject,
    backgroundColor: 'rgba(28,15,7,0.42)',
  },
  sheet: {
    backgroundColor: palette.cream100,
    borderTopLeftRadius: 24,
    borderTopRightRadius: 24,
    maxHeight: '92%',
    shadowColor: palette.espresso800,
    shadowOffset: { width: 0, height: -4 },
    shadowOpacity: 0.16,
    shadowRadius: 32,
    elevation: 12,
  },
  handle: {
    width: 36,
    height: 4,
    borderRadius: 9999,
    backgroundColor: palette.cream500,
    alignSelf: 'center',
    marginTop: 14,
    marginBottom: 4,
  },

  content: { paddingHorizontal: spacing[5], paddingBottom: 40, gap: 16 },

  title: {
    fontFamily: 'PlayfairDisplay_700Bold',
    fontSize: 24,
    color: colors.fgPrimary,
    marginTop: 8,
    marginBottom: 4,
  },
  hint: {
    fontFamily: 'DMSans_400Regular',
    fontSize: 12,
    color: colors.fgTertiary,
    fontStyle: 'italic',
    marginTop: -4,
  },

  fieldWrap: { gap: 6 },
  fieldLabel: {
    fontFamily: 'DMSans_500Medium',
    fontSize: 12,
    color: colors.fgSecondary,
  },
  required: { color: palette.caramel500 },
  input: {
    height: 50,
    paddingHorizontal: spacing[4],
    backgroundColor: palette.cream200,
    borderWidth: 1.5,
    borderColor: palette.cream400,
    borderRadius: 14,
    fontFamily: 'DMSans_400Regular',
    fontSize: 15,
    color: colors.fgPrimary,
  },
  inputFocused: {
    borderColor: palette.caramel400,
    shadowColor: palette.caramel400,
    shadowOffset: { width: 0, height: 0 },
    shadowOpacity: 0.12,
    shadowRadius: 3,
  },

  gearList: { gap: 8 },
  gearRow: {
    flexDirection: 'row',
    alignItems: 'center',
    backgroundColor: palette.cream200,
    borderWidth: 1.5,
    borderColor: palette.cream400,
    borderRadius: radii.lg,
    padding: 12,
    gap: 12,
  },
  gearRowSelected: {
    backgroundColor: palette.espresso800,
    borderColor: palette.espresso800,
  },
  gearBubble: {
    width: 44,
    height: 44,
    borderRadius: 11,
    backgroundColor: palette.cream300,
    alignItems: 'center',
    justifyContent: 'center',
  },
  gearBubbleSelected: {
    backgroundColor: palette.espresso700,
  },
  gearText: { flex: 1 },
  gearName: {
    fontFamily: 'DMSans_700Bold',
    fontSize: 14,
    color: colors.fgPrimary,
  },
  gearNameSelected: { color: palette.cream100 },
  gearSub: {
    fontFamily: 'DMSans_400Regular',
    fontSize: 12,
    color: colors.fgSecondary,
    marginTop: 1,
  },
  gearSubSelected: { color: palette.cream500 },
  toggle: {
    width: 20,
    height: 20,
    borderRadius: 10,
    borderWidth: 1.5,
    borderColor: palette.cream500,
    alignItems: 'center',
    justifyContent: 'center',
  },
  toggleSelected: {
    backgroundColor: palette.caramel400,
    borderColor: palette.caramel400,
  },
  toggleCheck: {
    fontSize: 11,
    color: '#fff',
    fontFamily: 'DMSans_700Bold',
  },

  errorBanner: {
    backgroundColor: palette.error100,
    borderRadius: radii.md,
    paddingVertical: 10,
    paddingHorizontal: 14,
  },
  errorText: {
    fontFamily: 'DMSans_500Medium',
    fontSize: 13,
    color: palette.error500,
  },

  cta: {
    height: 54,
    borderRadius: radii.xl,
    backgroundColor: palette.espresso800,
    alignItems: 'center',
    justifyContent: 'center',
  },
  ctaDisabled: { opacity: 0.38 },
  ctaPressed: { backgroundColor: palette.espresso700 },
  ctaLabel: {
    fontFamily: 'DMSans_700Bold',
    fontSize: 16,
    color: palette.cream100,
  },

  deleteBtn: {
    height: 46,
    borderRadius: radii.xl,
    backgroundColor: palette.error100,
    borderWidth: 1.5,
    borderColor: palette.error100,
    alignItems: 'center',
    justifyContent: 'center',
  },
  deleteBtnLabel: {
    fontFamily: 'DMSans_700Bold',
    fontSize: 15,
    color: palette.error500,
  },
});
