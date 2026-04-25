import {
  ActivityIndicator,
  Animated,
  KeyboardAvoidingView,
  LayoutAnimation,
  Platform,
  Pressable,
  ScrollView,
  StyleSheet,
  Text,
  TextInput,
  UIManager,
  View,
} from 'react-native';
import { useRef, useState } from 'react';
import { useSafeAreaInsets } from 'react-native-safe-area-context';
import { Feather } from '@expo/vector-icons';
import { colors, radii, spacing } from '@/src/theme';
import { useUser } from '@/src/context/UserContext';
import { authorizedFetch } from '@/src/lib/apiClient';
import type { User } from '@/src/context/UserContext';

if (Platform.OS === 'android') {
  UIManager.setLayoutAnimationEnabledExperimental?.(true);
}

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
      Animated.timing(opacity, { toValue: 1, duration: 260, useNativeDriver: true }),
      Animated.timing(translateY, { toValue: 0, duration: 260, useNativeDriver: true }),
    ]).start();
    timer.current = setTimeout(() => {
      Animated.parallel([
        Animated.timing(opacity, { toValue: 0, duration: 200, useNativeDriver: true }),
        Animated.timing(translateY, { toValue: 40, duration: 200, useNativeDriver: true }),
      ]).start();
    }, 2500);
  }

  const animatedStyle = { opacity, transform: [{ translateY }] };
  return { show, msg, animatedStyle };
}

// ─── Avatar ───────────────────────────────────────────────────────────────────

function Avatar({ user }: { user: User | null }) {
  const initials = user
    ? `${user.first_name.charAt(0)}${user.last_name.charAt(0)}`.toUpperCase()
    : '—';
  const fullName = user ? `${user.first_name} ${user.last_name}` : '';

  return (
    <View style={styles.avatarBlock}>
      <View style={styles.avatarCircle}>
        <Text style={styles.avatarInitials}>{initials}</Text>
      </View>
      <Text style={styles.avatarName}>{fullName}</Text>
      {user && <Text style={styles.avatarEmail}>{user.email}</Text>}
    </View>
  );
}

// ─── Section card shell ────────────────────────────────────────────────────────

interface SectionHeaderProps {
  label: string;
  preview: React.ReactNode;
  open: boolean;
  onToggle: () => void;
}

function SectionHeader({ label, preview, open, onToggle }: SectionHeaderProps) {
  return (
    <View style={styles.sectionRow}>
      <View style={styles.sectionLeft}>
        <Text style={styles.sectionLabel}>{label}</Text>
        {!open && preview}
      </View>
      <Pressable
        onPress={onToggle}
        style={({ pressed }) => [
          styles.sectionBtn,
          open ? styles.sectionBtnCancel : styles.sectionBtnEdit,
          pressed && { opacity: 0.75 },
        ]}
      >
        <Text style={[styles.sectionBtnText, open && styles.sectionBtnTextCancel]}>
          {open ? 'Cancel' : 'Edit'}
        </Text>
      </Pressable>
    </View>
  );
}

// ─── Screen ───────────────────────────────────────────────────────────────────

export default function ProfileScreen() {
  const { user, loading, setUser, logout } = useUser();
  const insets = useSafeAreaInsets();
  const toast = useToast();

  // Name section
  const [nameOpen, setNameOpen] = useState(false);
  const [nameForm, setNameForm] = useState({ firstName: '', lastName: '' });
  const [savingName, setSavingName] = useState(false);
  const [focusedName, setFocusedName] = useState<string | null>(null);
  const lastNameRef = useRef<TextInput>(null);

  // Password section
  const [pwOpen, setPwOpen] = useState(false);
  const [pwForm, setPwForm] = useState({ current: '', next: '', confirm: '' });
  const [pwError, setPwError] = useState<string | null>(null);
  const [savingPw, setSavingPw] = useState(false);
  const [focusedPw, setFocusedPw] = useState<string | null>(null);
  const [showCurrent, setShowCurrent] = useState(false);
  const [showNext, setShowNext] = useState(false);
  const [showConfirm, setShowConfirm] = useState(false);
  const nextPwRef = useRef<TextInput>(null);
  const confirmPwRef = useRef<TextInput>(null);

  const [loggingOut, setLoggingOut] = useState(false);

  function openName() {
    setNameForm({ firstName: user?.first_name ?? '', lastName: user?.last_name ?? '' });
    LayoutAnimation.configureNext(LayoutAnimation.Presets.easeInEaseOut);
    setNameOpen(true);
  }

  function closeNameAndReset() {
    LayoutAnimation.configureNext(LayoutAnimation.Presets.easeInEaseOut);
    setNameOpen(false);
  }

  function openPw() {
    setPwForm({ current: '', next: '', confirm: '' });
    setPwError(null);
    LayoutAnimation.configureNext(LayoutAnimation.Presets.easeInEaseOut);
    setPwOpen(true);
  }

  function closePwAndReset() {
    LayoutAnimation.configureNext(LayoutAnimation.Presets.easeInEaseOut);
    setPwOpen(false);
    setPwError(null);
  }

  async function saveName() {
    if (savingName || !nameForm.firstName.trim()) return;
    setSavingName(true);
    try {
      const updated = await authorizedFetch<User>('/user/me', {
        method: 'PATCH',
        body: { first_name: nameForm.firstName.trim(), last_name: nameForm.lastName.trim() },
      });
      setUser(updated);
      LayoutAnimation.configureNext(LayoutAnimation.Presets.easeInEaseOut);
      setNameOpen(false);
      toast.show('Name updated ✓');
    } catch (e) {
      toast.show(e instanceof Error ? e.message : 'Failed to update name.');
    } finally {
      setSavingName(false);
    }
  }

  async function savePassword() {
    if (savingPw) return;
    if (pwForm.next.length < 8) { setPwError('New password must be at least 8 characters.'); return; }
    if (pwForm.next !== pwForm.confirm) { setPwError("Passwords don't match."); return; }
    setPwError(null);
    setSavingPw(true);
    try {
      await authorizedFetch('/user/change-password', {
        method: 'POST',
        body: { old_password: pwForm.current, new_password: pwForm.next },
      });
      LayoutAnimation.configureNext(LayoutAnimation.Presets.easeInEaseOut);
      setPwOpen(false);
      setPwForm({ current: '', next: '', confirm: '' });
      toast.show('Password updated ✓');
    } catch (e) {
      setPwError(e instanceof Error ? e.message : 'Failed to update password.');
    } finally {
      setSavingPw(false);
    }
  }

  async function handleLogout() {
    if (loggingOut) return;
    setLoggingOut(true);
    await logout();
  }

  if (loading) {
    return (
      <View style={styles.center}>
        <ActivityIndicator color={colors.fgSecondary} />
      </View>
    );
  }

  const canSaveName = nameForm.firstName.trim().length > 0;
  const canSavePw = pwForm.current.length > 0 && pwForm.next.length > 0 && pwForm.confirm.length > 0;

  return (
    <View style={styles.screen}>
      <KeyboardAvoidingView
        style={{ flex: 1 }}
        behavior={Platform.OS === 'ios' ? 'padding' : 'height'}
      >
        <ScrollView
          contentContainerStyle={[styles.scroll, { paddingTop: insets.top + 12 }]}
          keyboardShouldPersistTaps="handled"
          showsVerticalScrollIndicator={false}
        >
          {/* Header */}
          <View style={styles.header}>
            <Text style={styles.title}>Profile</Text>
          </View>

          {/* Avatar */}
          <Avatar user={user} />

          {/* Name section */}
          <View style={styles.card}>
            <SectionHeader
              label="NAME"
              preview={
                <Text style={styles.previewText}>
                  {user ? `${user.first_name} ${user.last_name}` : '—'}
                </Text>
              }
              open={nameOpen}
              onToggle={nameOpen ? closeNameAndReset : openName}
            />

            {nameOpen && (
              <View style={styles.formWrap}>
                <View style={styles.nameRow}>
                  <View style={styles.nameCol}>
                    <TextInput
                      style={[styles.input, focusedName === 'firstName' && styles.inputFocused]}
                      placeholder="First name"
                      placeholderTextColor={colors.fgDisabled}
                      value={nameForm.firstName}
                      onChangeText={v => setNameForm(f => ({ ...f, firstName: v }))}
                      autoCapitalize="words"
                      returnKeyType="next"
                      onSubmitEditing={() => lastNameRef.current?.focus()}
                      onFocus={() => setFocusedName('firstName')}
                      onBlur={() => setFocusedName(null)}
                    />
                  </View>
                  <View style={styles.nameCol}>
                    <TextInput
                      ref={lastNameRef}
                      style={[styles.input, focusedName === 'lastName' && styles.inputFocused]}
                      placeholder="Last name"
                      placeholderTextColor={colors.fgDisabled}
                      value={nameForm.lastName}
                      onChangeText={v => setNameForm(f => ({ ...f, lastName: v }))}
                      autoCapitalize="words"
                      returnKeyType="done"
                      onSubmitEditing={saveName}
                      onFocus={() => setFocusedName('lastName')}
                      onBlur={() => setFocusedName(null)}
                    />
                  </View>
                </View>

                <Pressable
                  style={({ pressed }) => [
                    styles.cta,
                    (!canSaveName || savingName) && styles.ctaDisabled,
                    pressed && canSaveName && !savingName && styles.ctaPressed,
                  ]}
                  onPress={saveName}
                  disabled={!canSaveName || savingName}
                >
                  {savingName
                    ? <ActivityIndicator color={colors.fgInverse} />
                    : <Text style={styles.ctaLabel}>Save name</Text>}
                </Pressable>
              </View>
            )}
          </View>

          {/* Password section */}
          <View style={styles.card}>
            <SectionHeader
              label="PASSWORD"
              preview={
                <Text style={[styles.previewText, styles.passwordDots]}>••••••••</Text>
              }
              open={pwOpen}
              onToggle={pwOpen ? closePwAndReset : openPw}
            />

            {pwOpen && (
              <View style={styles.formWrap}>
                {/* Current password */}
                <View style={styles.passwordInputWrap}>
                  <TextInput
                    style={[styles.input, styles.passwordInput, focusedPw === 'current' && styles.inputFocused]}
                    placeholder="Your current password"
                    placeholderTextColor={colors.fgDisabled}
                    value={pwForm.current}
                    onChangeText={v => { setPwForm(f => ({ ...f, current: v })); setPwError(null); }}
                    secureTextEntry={!showCurrent}
                    returnKeyType="next"
                    onSubmitEditing={() => nextPwRef.current?.focus()}
                    onFocus={() => setFocusedPw('current')}
                    onBlur={() => setFocusedPw(null)}
                  />
                  <Pressable style={styles.eyeBtn} onPress={() => setShowCurrent(v => !v)}>
                    <Feather name={showCurrent ? 'eye-off' : 'eye'} size={18} color={colors.fgTertiary} />
                  </Pressable>
                </View>

                {/* New password */}
                <View style={styles.passwordInputWrap}>
                  <TextInput
                    ref={nextPwRef}
                    style={[styles.input, styles.passwordInput, focusedPw === 'next' && styles.inputFocused]}
                    placeholder="Min. 8 characters"
                    placeholderTextColor={colors.fgDisabled}
                    value={pwForm.next}
                    onChangeText={v => { setPwForm(f => ({ ...f, next: v })); setPwError(null); }}
                    secureTextEntry={!showNext}
                    returnKeyType="next"
                    onSubmitEditing={() => confirmPwRef.current?.focus()}
                    onFocus={() => setFocusedPw('next')}
                    onBlur={() => setFocusedPw(null)}
                  />
                  <Pressable style={styles.eyeBtn} onPress={() => setShowNext(v => !v)}>
                    <Feather name={showNext ? 'eye-off' : 'eye'} size={18} color={colors.fgTertiary} />
                  </Pressable>
                </View>

                {/* Confirm password */}
                <View style={styles.passwordInputWrap}>
                  <TextInput
                    ref={confirmPwRef}
                    style={[styles.input, styles.passwordInput, focusedPw === 'confirm' && styles.inputFocused]}
                    placeholder="Same again"
                    placeholderTextColor={colors.fgDisabled}
                    value={pwForm.confirm}
                    onChangeText={v => { setPwForm(f => ({ ...f, confirm: v })); setPwError(null); }}
                    secureTextEntry={!showConfirm}
                    returnKeyType="done"
                    onSubmitEditing={savePassword}
                    onFocus={() => setFocusedPw('confirm')}
                    onBlur={() => setFocusedPw(null)}
                  />
                  <Pressable style={styles.eyeBtn} onPress={() => setShowConfirm(v => !v)}>
                    <Feather name={showConfirm ? 'eye-off' : 'eye'} size={18} color={colors.fgTertiary} />
                  </Pressable>
                </View>

                {pwError && (
                  <View style={styles.errorBanner}>
                    <Text style={styles.errorText}>{pwError}</Text>
                  </View>
                )}

                <Pressable
                  style={({ pressed }) => [
                    styles.cta,
                    (!canSavePw || savingPw) && styles.ctaDisabled,
                    pressed && canSavePw && !savingPw && styles.ctaPressed,
                  ]}
                  onPress={savePassword}
                  disabled={!canSavePw || savingPw}
                >
                  {savingPw
                    ? <ActivityIndicator color={colors.fgInverse} />
                    : <Text style={styles.ctaLabel}>Update password</Text>}
                </Pressable>
              </View>
            )}
          </View>

          {/* Log out */}
          <Pressable
            style={({ pressed }) => [
              styles.logoutBtn,
              pressed && styles.logoutPressed,
              loggingOut && styles.ctaDisabled,
            ]}
            onPress={handleLogout}
            disabled={loggingOut}
          >
            {loggingOut
              ? <ActivityIndicator color="#C0392B" />
              : <Text style={styles.logoutLabel}>Log out</Text>}
          </Pressable>
        </ScrollView>
      </KeyboardAvoidingView>

      {/* Toast */}
      <Animated.View style={[styles.toast, toast.animatedStyle]} pointerEvents="none">
        <Text style={styles.toastText}>{toast.msg}</Text>
      </Animated.View>
    </View>
  );
}

// ─── Styles ───────────────────────────────────────────────────────────────────

const styles = StyleSheet.create({
  screen: {
    flex: 1,
    backgroundColor: colors.bgApp,
  },
  center: {
    flex: 1,
    backgroundColor: colors.bgApp,
    alignItems: 'center',
    justifyContent: 'center',
  },
  scroll: {
    paddingBottom: 32,
  },

  // Header
  header: {
    paddingTop: 12,
    paddingHorizontal: spacing[6],
    paddingBottom: spacing[4],
  },
  title: {
    fontFamily: 'PlayfairDisplay_900Black',
    fontSize: 34,
    lineHeight: 40,
    letterSpacing: -0.7,
    color: colors.fgPrimary,
  },

  // Avatar
  avatarBlock: {
    alignItems: 'center',
    paddingHorizontal: spacing[5],
    paddingBottom: spacing[6],
    paddingTop: spacing[2],
  },
  avatarCircle: {
    width: 72,
    height: 72,
    borderRadius: 36,
    backgroundColor: colors.fgPrimary,
    alignItems: 'center',
    justifyContent: 'center',
    shadowColor: colors.fgPrimary,
    shadowOffset: { width: 0, height: 4 },
    shadowOpacity: 0.2,
    shadowRadius: 16,
    elevation: 6,
  },
  avatarInitials: {
    fontFamily: 'DMSans_700Bold',
    fontSize: 26,
    color: colors.fgInverse,
  },
  avatarName: {
    fontFamily: 'DMSans_700Bold',
    fontSize: 18,
    color: colors.fgPrimary,
    marginTop: 14,
  },
  avatarEmail: {
    fontFamily: 'DMSans_400Regular',
    fontSize: 13,
    color: colors.fgSecondary,
    marginTop: 4,
  },

  // Section card
  card: {
    marginHorizontal: spacing[5],
    marginBottom: 12,
    backgroundColor: colors.bgCard,
    borderRadius: radii.lg,
    padding: 18,
    shadowColor: colors.fgPrimary,
    shadowOffset: { width: 0, height: 2 },
    shadowOpacity: 0.07,
    shadowRadius: 8,
    elevation: 2,
  },
  sectionRow: {
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'space-between',
  },
  sectionLeft: {
    flex: 1,
    marginRight: 12,
  },
  sectionLabel: {
    fontFamily: 'DMSans_700Bold',
    fontSize: 10,
    color: colors.fgSecondary,
    letterSpacing: 0.08 * 10,
    textTransform: 'uppercase',
    marginBottom: 3,
  },
  previewText: {
    fontFamily: 'DMSans_700Bold',
    fontSize: 15,
    color: colors.fgPrimary,
  },
  passwordDots: {
    letterSpacing: 2,
  },
  sectionBtn: {
    height: 30,
    paddingHorizontal: 14,
    borderRadius: radii.full,
    alignItems: 'center',
    justifyContent: 'center',
  },
  sectionBtnEdit: {
    backgroundColor: colors.bgSubtle,
  },
  sectionBtnCancel: {
    backgroundColor: 'transparent',
    borderWidth: 1.5,
    borderColor: colors.borderSubtle,
  },
  sectionBtnText: {
    fontFamily: 'DMSans_700Bold',
    fontSize: 12,
    color: colors.fgPrimary,
  },
  sectionBtnTextCancel: {
    color: colors.fgSecondary,
  },

  // Form
  formWrap: {
    marginTop: 18,
    gap: 12,
  },
  nameRow: {
    flexDirection: 'row',
    gap: 12,
  },
  nameCol: {
    flex: 1,
  },
  input: {
    height: 50,
    paddingHorizontal: spacing[4],
    backgroundColor: colors.bgCard,
    borderWidth: 1.5,
    borderColor: colors.borderSubtle,
    borderRadius: 14,
    fontFamily: 'DMSans_400Regular',
    fontSize: 15,
    color: colors.fgPrimary,
  },
  inputFocused: {
    borderColor: colors.borderFocus,
    shadowColor: colors.borderFocus,
    shadowOffset: { width: 0, height: 0 },
    shadowOpacity: 0.12,
    shadowRadius: 3,
    elevation: 0,
  },
  passwordInputWrap: {
    position: 'relative',
  },
  passwordInput: {
    paddingRight: 48,
  },
  eyeBtn: {
    position: 'absolute',
    right: 14,
    top: 0,
    bottom: 0,
    justifyContent: 'center',
    paddingHorizontal: 2,
  },

  // Error
  errorBanner: {
    backgroundColor: '#FDECEA',
    borderWidth: 1,
    borderColor: '#F5C6C2',
    borderRadius: radii.md,
    paddingVertical: 10,
    paddingHorizontal: 14,
  },
  errorText: {
    fontFamily: 'DMSans_500Medium',
    fontSize: 13,
    color: '#C0392B',
  },

  // CTA
  cta: {
    height: 50,
    borderRadius: radii.full,
    backgroundColor: colors.interactivePrimary,
    alignItems: 'center',
    justifyContent: 'center',
  },
  ctaDisabled: {
    opacity: 0.38,
  },
  ctaPressed: {
    transform: [{ scale: 0.97 }],
    backgroundColor: colors.interactivePrimaryHover,
  },
  ctaLabel: {
    fontFamily: 'DMSans_700Bold',
    fontSize: 15,
    color: colors.fgInverse,
  },

  // Logout
  logoutBtn: {
    height: 50,
    borderRadius: radii.full,
    backgroundColor: '#FDECEA',
    borderWidth: 1.5,
    borderColor: '#FDECEA',
    alignItems: 'center',
    justifyContent: 'center',
    marginHorizontal: spacing[5],
    marginTop: 8,
    marginBottom: 32,
  },
  logoutPressed: {
    backgroundColor: '#f5c6c2',
    borderColor: '#f5c6c2',
  },
  logoutLabel: {
    fontFamily: 'DMSans_700Bold',
    fontSize: 15,
    color: '#C0392B',
  },

  // Toast
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
