package com.triforceblitz.triforceblitz.randomizer;

import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;

import java.util.Map;

import static org.assertj.core.api.Assertions.assertThat;
import static org.assertj.core.api.Assertions.assertThatThrownBy;

class PresetTests {
    Preset preset;

    @BeforeEach
    void setUp() {
        preset = new Preset(Map.of(
                "create_spoiler", true,
                "user_message", "Test Settings",
                "bridge_medallions", 6
        ));
    }

    @Test
    void getSetting_whenExists_returnsSetting() {
        assertThat(preset.getSetting("bridge_medallions", Integer.class)).isPresent();
    }

    @Test
    void getSetting_whenMissing_returnsEmpty() {
        assertThat(preset.getSetting("invalid", Object.class)).isEmpty();
    }

    @Test
    void getSetting_whenInvalidType_throwsException() {
        assertThatThrownBy(() -> preset.getSetting("create_spoiler", Integer.class))
                .isInstanceOf(ClassCastException.class);
    }

    @Test
    void isMultiWorld_whenWorldCountIs2_returnsTrue() {
        preset.setSetting("world_count", 2);
        assertThat(preset.isMultiWorld()).isTrue();
    }

    @Test
    void isMultiWorld_whenWorldCountIsOne_returnsFalse() {
        preset.setSetting("world_count", 1);
        assertThat(preset.isMultiWorld()).isFalse();
    }

    @Test
    void isMultiWorld_whenWorldCountIsNotInteger_returnsFalse() {
        preset.setSetting("world_count", "many");
        assertThat(preset.isMultiWorld()).isFalse();
    }

    @Test
    void enable_enablesPreset() {
        preset.enable();
        assertThat(preset.isEnabled()).isTrue();
        assertThat(preset.isDisabled()).isFalse();
    }

    @Test
    void disable_disablesPreset() {
        preset.disable();
        assertThat(preset.isEnabled()).isFalse();
        assertThat(preset.isDisabled()).isTrue();
    }

    @Test
    void compareTo_whenLessThan_isNegative() {
        preset.setOrdinal(100);
        var other = new Preset();
        other.setOrdinal(200);
        assertThat(preset.compareTo(other)).isNegative();
    }

    @Test
    void compareTo_whenGreaterThan_isPositive() {
        preset.setOrdinal(200);
        var other = new Preset();
        other.setOrdinal(100);
        assertThat(preset.compareTo(other)).isPositive();
    }
}