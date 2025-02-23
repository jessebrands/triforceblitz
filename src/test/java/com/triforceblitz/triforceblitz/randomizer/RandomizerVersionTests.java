package com.triforceblitz.triforceblitz.randomizer;

import org.junit.jupiter.api.Test;

import static org.assertj.core.api.Assertions.assertThat;
import static org.assertj.core.api.Assertions.assertThatThrownBy;

class RandomizerVersionTests {
    @Test
    void of_whenValid_createsVersion() {
        var version = RandomizerVersion.of("1.2.3-blitz-4.5");
        var expected = new RandomizerVersion(1, 2, 3, "blitz", 4, 5);
        assertThat(version).isEqualTo(expected);
    }

    @Test
    void of_whenInvalid_throwsIllegalArgumentException() {
        assertThatThrownBy(() -> RandomizerVersion.of("1.2.3 blitz-45"))
                .isInstanceOf(IllegalArgumentException.class);
    }

    @Test
    void valid_whenValid_returnsTrue() {
        assertThat(RandomizerVersion.valid("6.2.0-blitz-0.12")).isTrue();
    }

    @Test
    void valid_whenInvalid_returnsFalse() {
        assertThat(RandomizerVersion.valid("8.2.42 blitz-22")).isFalse();
    }

    @Test
    void toString_isCorrectFormat() {
        var version = new RandomizerVersion(8, 2, 12, "blitz", 0, 60);
        assertThat(version.toString()).isEqualTo("8.2.12-blitz-0.60");
    }
}