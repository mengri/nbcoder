package notify

import "testing"

func TestSubscriptionPreference_MuteChannel(t *testing.T) {
	pref := NewSubscriptionPreference("pref-1", "user-1", "CardCreated")
	pref.MuteChannel(ChannelEmail)
	if !pref.IsChannelMuted(ChannelEmail) {
		t.Error("expected email channel to be muted")
	}
	if pref.IsChannelMuted(ChannelSystem) {
		t.Error("expected system channel to not be muted")
	}
}

func TestSubscriptionPreference_UnmuteChannel(t *testing.T) {
	pref := NewSubscriptionPreference("pref-1", "user-1", "CardCreated")
	pref.MuteChannel(ChannelEmail)
	pref.MuteChannel(ChannelSystem)
	pref.UnmuteChannel(ChannelEmail)
	if pref.IsChannelMuted(ChannelEmail) {
		t.Error("expected email channel to be unmuted")
	}
	if !pref.IsChannelMuted(ChannelSystem) {
		t.Error("expected system channel to still be muted")
	}
}

func TestSubscriptionPreference_MuteChannel_Idempotent(t *testing.T) {
	pref := NewSubscriptionPreference("pref-1", "user-1", "CardCreated")
	pref.MuteChannel(ChannelEmail)
	pref.MuteChannel(ChannelEmail)
	if len(pref.MutedChannels) != 1 {
		t.Errorf("expected 1 muted channel, got %d", len(pref.MutedChannels))
	}
}

func TestSubscription_UnmuteChannel_NotMuted(t *testing.T) {
	pref := NewSubscriptionPreference("pref-1", "user-1", "CardCreated")
	pref.UnmuteChannel(ChannelEmail)
	if len(pref.MutedChannels) != 0 {
		t.Errorf("expected 0 muted channels, got %d", len(pref.MutedChannels))
	}
}

func TestSubscription_IsMutedForEvent(t *testing.T) {
	sub := NewSubscription("sub-1", "user-1", "CardCreated", ChannelEmail)
	if sub.IsMutedForEvent("CardCreated") {
		t.Error("expected subscription to not be muted")
	}
	sub.Mute()
	if !sub.IsMutedForEvent("CardCreated") {
		t.Error("expected subscription to be muted after Mute()")
	}
}
