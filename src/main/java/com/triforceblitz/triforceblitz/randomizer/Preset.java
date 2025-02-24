package com.triforceblitz.triforceblitz.randomizer;

import java.util.HashMap;
import java.util.Map;
import java.util.Optional;

/**
 * Settings preset for the Ocarina of Time Randomizer.
 */
public class Preset implements Comparable<Preset> {
    private static final String WORLD_COUNT_SETTING_NAME = "world_count";

    /// Randomizer settings map.
    private final Map<String, Object> settings = new HashMap<>();

    /// Flag to set whether the preset can be used.
    private boolean enabled;

    /// Ordinal for sorting presets.
    private int ordinal;

    /**
     * Constructs an empty preset.
     */
    public Preset() {
    }

    /**
     * Constructs a preset with settings.
     *
     * @param settings Randomizer settings.
     */
    public Preset(Map<String, Object> settings) {
        this.settings.putAll(settings);
    }

    /**
     * Gets a randomizer setting from the preset.
     *
     * @param name Name of the setting.
     * @return {@link Optional} containing the value, or empty if not set.
     */
    public Optional<Object> getSetting(String name) {
        return Optional.ofNullable(settings.get(name));
    }

    /**
     * Gets a randomizer setting from the preset.
     *
     * @param name The name of the setting.
     * @param cls  Class type reference to cast to.
     * @return Setting value or empty if not set.
     * @throws ClassCastException if present but different type.
     */
    public <T> Optional<T> getSetting(String name, Class<T> cls) {
        var value = settings.get(name);
        if (value == null) {
            return Optional.empty();
        } else if (cls.isInstance(value)) {
            return Optional.of(cls.cast(value));
        }
        throw new ClassCastException();
    }

    /**
     * Sets the value of a randomizer setting.
     *
     * @param name  Name of the setting.
     * @param value Value to set the setting to.
     */
    public void setSetting(String name, Object value) {
        settings.put(name, value);
    }

    /**
     * Checks if this is a multi-world preset.
     *
     * @return <code>true</code> if multi-world, <code>false</code> if solo.
     */
    public boolean isMultiWorld() {
        return getSetting(WORLD_COUNT_SETTING_NAME)
                .filter(s -> s instanceof Integer count && count >= 2)
                .isPresent();
    }

    /**
     * Checks whether the preset is enabled for use.
     *
     * @return <code>true</code> if enabled, <code>false</code> if disabled.
     */
    public boolean isEnabled() {
        return enabled;
    }

    /**
     * Checks whether the preset is disabled for use.
     *
     * @return <code>true</code> if disabled, <code>false</code> if enabled.
     */
    public boolean isDisabled() {
        return !enabled;
    }

    /**
     * Enables the preset.
     */
    public void enable() {
        this.enabled = true;
    }

    /**
     * Disables the preset.
     */
    public void disable() {
        this.enabled = false;
    }

    /**
     * Sets the sorting ordinal.
     *
     * @param ordinal Ordinal value used for sorting.
     */
    public void setOrdinal(int ordinal) {
        this.ordinal = ordinal;
    }

    @Override
    public int compareTo(Preset other) {
        return Integer.compare(ordinal, other.ordinal);
    }
}
