package com.triforceblitz.triforceblitz.randomizer;

import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;

import static org.assertj.core.api.Assertions.assertThat;

class RandomizerTests {
    private Randomizer randomizer;

    @BeforeEach
    void setUp() {
        randomizer = new Randomizer(RandomizerVersion.of("1.0.0-blitz-1.0"));
    }

    @Test
    void getVersion_returnsVersion() {
        var version = randomizer.getVersion();
        var expected = new RandomizerVersion(1, 0, 0, "blitz", 1, 0);
        assertThat(version).isEqualTo(expected);
    }

    @Test
    void enable_enablesRandomizer() {
        randomizer.enable();
        assertThat(randomizer.isEnabled()).isTrue();
    }

    @Test
    void disable_disablesRandomizer() {
        randomizer.disable();
        assertThat(randomizer.isEnabled()).isFalse();
    }

    @Test
    void setPrerelease_setsPrereleaseStatus() {
        randomizer.setPrerelease(true);
        assertThat(randomizer.isPrerelease()).isTrue();
    }

    @Test
    void equals_whenEqual_returnsTrue() {
        var other = new Randomizer(RandomizerVersion.of("1.0.0-blitz-1.0"));
        assertThat(randomizer).isEqualTo(other);
    }

    @Test
    void equals_whenNotEqual_returnsFalse() {
        var string = "1.0.0-blitz-1.0";
        var other =  new Randomizer(RandomizerVersion.of("1.2.1-blitz-1.0"));
        assertThat(randomizer).isNotEqualTo(string);
        assertThat(randomizer).isNotEqualTo(other);
    }

    @Test
    void hashCode_whenEqual_isEqual() {
        var other = new Randomizer(RandomizerVersion.of("1.0.0-blitz-1.0"));
        assertThat(randomizer.hashCode()).isEqualTo(other.hashCode());
    }

    @Test
    void hashCode_whenNotEqual_isNotEqual() {
        var other =  new Randomizer(RandomizerVersion.of("1.2.1-blitz-1.0"));
        assertThat(randomizer.hashCode()).isNotEqualTo(other.hashCode());
    }

    @Test
    void toString_isCorrectFormat() {
        var version = RandomizerVersion.of("1.0.0-blitz-1.0");
        var expected = version.toString();
        assertThat(randomizer.toString()).isEqualTo(expected);
    }
}