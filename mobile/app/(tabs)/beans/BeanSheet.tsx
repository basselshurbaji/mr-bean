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
import { beansApi, Bean, BeanBody, PROCESSES, ROAST_LEVELS } from '@/src/api/beans';

interface Props {
  editBean?: Bean;
  onClose: () => void;
  onSaved: (bean: Bean) => void;
}

export default function BeanSheet({ editBean, onClose, onSaved }: Props) {
  const translateY      = useRef(new Animated.Value(700)).current;
  const backdropOpacity = useRef(new Animated.Value(0)).current;

  const isEdit = !!editBean;

  const [form, setForm] = useState({
    name:         editBean?.name          ?? '',
    roaster:      editBean?.roaster       ?? '',
    origin:       editBean?.origin        ?? '',
    tastingNotes: editBean?.tasting_notes ?? '',
    notes:        editBean?.notes         ?? '',
  });
  const [selectedProcess,    setSelectedProcess]    = useState<string | null>(editBean?.process    ?? null);
  const [selectedRoastLevel, setSelectedRoastLevel] = useState<string | null>(editBean?.roast_level ?? null);
  const [focused, setFocused]  = useState<string | null>(null);
  const [saving,  setSaving]   = useState(false);
  const [error,   setError]    = useState<string | null>(null);

  const nameRef         = useRef<TextInput>(null);
  const roasterRef      = useRef<TextInput>(null);
  const originRef       = useRef<TextInput>(null);
  const tastingNotesRef = useRef<TextInput>(null);
  const notesRef        = useRef<TextInput>(null);

  useEffect(() => {
    Animated.parallel([
      Animated.timing(translateY,      { toValue: 0, duration: 300, useNativeDriver: true }),
      Animated.timing(backdropOpacity, { toValue: 1, duration: 300, useNativeDriver: true }),
    ]).start();
  }, []);

  function dismiss() {
    Animated.parallel([
      Animated.timing(translateY,      { toValue: 700, duration: 260, useNativeDriver: true }),
      Animated.timing(backdropOpacity, { toValue: 0,   duration: 260, useNativeDriver: true }),
    ]).start(() => onClose());
  }

  async function save() {
    if (saving || !form.name.trim()) return;
    setSaving(true);
    setError(null);
    try {
      const body: BeanBody = {
        name:         form.name.trim(),
        roaster:      form.roaster.trim()      || undefined,
        origin:       form.origin.trim()       || undefined,
        process:      selectedProcess          || undefined,
        roast_level:  selectedRoastLevel       || undefined,
        tasting_notes: form.tastingNotes.trim() || undefined,
        notes:        form.notes.trim()        || undefined,
      };
      const result = isEdit
        ? await beansApi.update(editBean!.id, body)
        : await beansApi.create(body);
      onSaved(result);
    } catch (e) {
      setError(e instanceof Error ? e.message : 'Something went wrong.');
      setSaving(false);
    }
  }

  const canSave = form.name.trim().length > 0;

  return (
    <Modal transparent animationType="none" onRequestClose={dismiss}>
      <View style={styles.overlay}>
        <Animated.View style={[styles.backdrop, { opacity: backdropOpacity }]}>
          <Pressable style={StyleSheet.absoluteFill} onPress={dismiss} />
        </Animated.View>

        <Animated.View style={[styles.sheet, { transform: [{ translateY }] }]}>
          <View style={styles.handle} />

          <KeyboardAvoidingView behavior={Platform.OS === 'ios' ? 'padding' : 'height'}>
            <ScrollView
              contentContainerStyle={styles.formContent}
              showsVerticalScrollIndicator={false}
              keyboardShouldPersistTaps="handled"
            >
              {/* Header row */}
              <View style={styles.sheetHeader}>
                <Text style={styles.sheetTitle}>
                  {isEdit ? 'Edit bean' : 'Add a bean'}
                </Text>
                <Pressable
                  style={({ pressed }) => [styles.closeBtn, pressed && { opacity: 0.7 }]}
                  onPress={dismiss}
                >
                  <Text style={styles.closeBtnIcon}>✕</Text>
                </Pressable>
              </View>

              {/* Name */}
              <View style={styles.fieldWrap}>
                <Text style={styles.fieldLabel}>
                  Name <Text style={styles.required}>*</Text>
                </Text>
                <TextInput
                  ref={nameRef}
                  style={[styles.input, focused === 'name' && styles.inputFocused]}
                  placeholder="e.g. Ethiopia Yirgacheffe"
                  placeholderTextColor={palette.cream600}
                  value={form.name}
                  onChangeText={v => setForm(f => ({ ...f, name: v }))}
                  onFocus={() => setFocused('name')}
                  onBlur={() => setFocused(null)}
                  returnKeyType="next"
                  onSubmitEditing={() => roasterRef.current?.focus()}
                  textContentType="none"
                  autoComplete="off"
                />
              </View>

              {/* Roaster + Origin */}
              <View style={styles.row}>
                <View style={styles.rowCol}>
                  <Text style={styles.fieldLabel}>Roaster</Text>
                  <TextInput
                    ref={roasterRef}
                    style={[styles.input, focused === 'roaster' && styles.inputFocused]}
                    placeholder="e.g. Onyx"
                    placeholderTextColor={palette.cream600}
                    value={form.roaster}
                    onChangeText={v => setForm(f => ({ ...f, roaster: v }))}
                    onFocus={() => setFocused('roaster')}
                    onBlur={() => setFocused(null)}
                    returnKeyType="next"
                    onSubmitEditing={() => originRef.current?.focus()}
                    textContentType="none"
                    autoComplete="off"
                  />
                </View>
                <View style={styles.rowCol}>
                  <Text style={styles.fieldLabel}>Origin</Text>
                  <TextInput
                    ref={originRef}
                    style={[styles.input, focused === 'origin' && styles.inputFocused]}
                    placeholder="e.g. Ethiopia"
                    placeholderTextColor={palette.cream600}
                    value={form.origin}
                    onChangeText={v => setForm(f => ({ ...f, origin: v }))}
                    onFocus={() => setFocused('origin')}
                    onBlur={() => setFocused(null)}
                    returnKeyType="next"
                    onSubmitEditing={() => tastingNotesRef.current?.focus()}
                    textContentType="none"
                    autoComplete="off"
                  />
                </View>
              </View>

              {/* Process chips */}
              <View style={styles.fieldWrap}>
                <Text style={styles.fieldLabel}>Process</Text>
                <View style={styles.chips}>
                  {PROCESSES.map(p => {
                    const active = selectedProcess === p.id;
                    return (
                      <Pressable
                        key={p.id}
                        style={[styles.chip, active && styles.chipActive]}
                        onPress={() => setSelectedProcess(active ? null : p.id)}
                      >
                        <Text style={[styles.chipLabel, active && styles.chipLabelActive]}>
                          {p.label}
                        </Text>
                      </Pressable>
                    );
                  })}
                </View>
              </View>

              {/* Roast level chips */}
              <View style={styles.fieldWrap}>
                <Text style={styles.fieldLabel}>Roast level</Text>
                <View style={styles.chips}>
                  {ROAST_LEVELS.map(r => {
                    const active = selectedRoastLevel === r.id;
                    return (
                      <Pressable
                        key={r.id}
                        style={[styles.chip, active && styles.chipActive]}
                        onPress={() => setSelectedRoastLevel(active ? null : r.id)}
                      >
                        <View style={[
                          styles.roastDot,
                          { backgroundColor: active ? palette.cream100 : r.color },
                        ]} />
                        <Text style={[styles.chipLabel, active && styles.chipLabelActive]}>
                          {r.label}
                        </Text>
                      </Pressable>
                    );
                  })}
                </View>
              </View>

              {/* Tasting notes */}
              <View style={styles.fieldWrap}>
                <Text style={styles.fieldLabel}>Tasting notes</Text>
                <TextInput
                  ref={tastingNotesRef}
                  style={[styles.input, styles.textarea, focused === 'tasting' && styles.inputFocused]}
                  placeholder="e.g. Jasmine, bergamot, peach."
                  placeholderTextColor={palette.cream600}
                  value={form.tastingNotes}
                  onChangeText={v => setForm(f => ({ ...f, tastingNotes: v }))}
                  onFocus={() => setFocused('tasting')}
                  onBlur={() => setFocused(null)}
                  multiline
                  numberOfLines={2}
                  textAlignVertical="top"
                  returnKeyType="next"
                  onSubmitEditing={() => notesRef.current?.focus()}
                  textContentType="none"
                  autoComplete="off"
                />
              </View>

              {/* Notes */}
              <View style={styles.fieldWrap}>
                <Text style={styles.fieldLabel}>Notes</Text>
                <TextInput
                  ref={notesRef}
                  style={[styles.input, styles.textarea, focused === 'notes' && styles.inputFocused]}
                  placeholder="Your impressions, ratios, anything worth remembering."
                  placeholderTextColor={palette.cream600}
                  value={form.notes}
                  onChangeText={v => setForm(f => ({ ...f, notes: v }))}
                  onFocus={() => setFocused('notes')}
                  onBlur={() => setFocused(null)}
                  multiline
                  numberOfLines={2}
                  textAlignVertical="top"
                  textContentType="none"
                  autoComplete="off"
                />
              </View>

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
                    {isEdit ? 'Save changes' : 'Add to your beans'}
                  </Text>
                )}
              </Pressable>
            </ScrollView>
          </KeyboardAvoidingView>
        </Animated.View>
      </View>
    </Modal>
  );
}

const styles = StyleSheet.create({
  overlay:  { flex: 1, justifyContent: 'flex-end' },
  backdrop: {
    ...StyleSheet.absoluteFillObject,
    backgroundColor: 'rgba(28,15,7,0.42)',
  },
  sheet: {
    backgroundColor: palette.cream100,
    borderTopLeftRadius: 28,
    borderTopRightRadius: 28,
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
    borderRadius: radii.full,
    backgroundColor: palette.cream500,
    alignSelf: 'center',
    marginTop: 14,
    marginBottom: 4,
  },

  formContent: { paddingHorizontal: spacing[5], paddingBottom: 40, gap: 16 },

  sheetHeader: {
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'space-between',
    marginBottom: 4,
    marginTop: 8,
  },
  sheetTitle: {
    fontFamily: 'PlayfairDisplay_700Bold',
    fontSize: 24,
    color: colors.fgPrimary,
  },
  closeBtn: {
    width: 34,
    height: 34,
    borderRadius: 17,
    backgroundColor: palette.cream300,
    alignItems: 'center',
    justifyContent: 'center',
  },
  closeBtnIcon: {
    fontFamily: 'DMSans_400Regular',
    fontSize: 14,
    color: palette.espresso500,
  },

  fieldWrap: { gap: 8 },
  fieldLabel: {
    fontFamily: 'DMSans_500Medium',
    fontSize: 13,
    color: colors.fgSecondary,
  },
  required: { color: palette.caramel500 },

  row:    { flexDirection: 'row', gap: 12 },
  rowCol: { flex: 1, gap: 8 },

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
  textarea: {
    height: 80,
    paddingTop: 14,
    paddingBottom: 14,
  },
  inputFocused: {
    borderColor: palette.caramel400,
    shadowColor: palette.caramel400,
    shadowOffset: { width: 0, height: 0 },
    shadowOpacity: 0.12,
    shadowRadius: 3,
  },

  chips: {
    flexDirection: 'row',
    flexWrap: 'wrap',
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
    flexDirection: 'row',
    gap: 6,
  },
  chipActive: {
    backgroundColor: palette.espresso800,
    borderColor: palette.espresso800,
  },
  chipLabel: {
    fontFamily: 'DMSans_700Bold',
    fontSize: 13,
    color: palette.espresso500,
  },
  chipLabelActive: { color: palette.cream100 },

  roastDot: {
    width: 10,
    height: 10,
    borderRadius: 5,
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
    marginTop: 4,
  },
  ctaDisabled: { opacity: 0.38 },
  ctaPressed:  { backgroundColor: palette.espresso700 },
  ctaLabel: {
    fontFamily: 'DMSans_700Bold',
    fontSize: 16,
    color: palette.cream100,
  },
});
