package com.triforceblitz.triforceblitz.randomizer;

import java.util.Objects;

/**
 * Core domain object of the randomizer package.
 *
 * <p>Provides an abstraction around an Ocarina of Time Randomizer that can be
 * invoked to generate a new ROM. Randomizers also have setting presets that
 * can be queried and modified, allowing one to quickly generate a ROM or patch
 * file from a template.</p>
 *
 * <p>Randomizers must be implemented from the outside as they invoke an
 * external process, this is left as an implementation detail.</p>
 */
public class Randomizer {
    /// Version that uniquely identifies the randomizer.
    private final RandomizerVersion version;

    /// Flag that sets whether the randomizer can be used to generate ROMs.
    private boolean enabled;

    /// Flag that sets whether the randomizer is a prerelease version.
    private boolean prerelease;

    public Randomizer(RandomizerVersion version) {
        this.version = version;
    }

    /**
     * Gets the version of the randomizer.
     *
     * @return The randomizer version.
     */
    public RandomizerVersion getVersion() {
        return version;
    }

    /**
     * Returns whether this Randomizer is allowed to generate ROMs or not.
     *
     * @return <code>true</code> if enabled, <code>false</code> if not.
     */
    public boolean isEnabled() {
        return enabled;
    }

    /**
     * Enables the Randomizer.
     */
    public void enable() {
        this.enabled = true;
    }

    /**
     * Disables the Randomizer.
     */
    public void disable() {
        this.enabled = false;
    }

    /**
     * Returns the Randomizer's prerelease status.
     *
     * @return <code>true</code> if prerelease version, <code>false</code> if
     * not.
     */
    public boolean isPrerelease() {
        return prerelease;
    }

    /**
     * Sets the Randomizer's prerelease status.
     *
     * @param prerelease <code>true</code> to flag the Randomizer as a
     *                   prerelease version, <code>false</code> if not.
     */
    public void setPrerelease(boolean prerelease) {
        this.prerelease = prerelease;
    }

    @Override
    public boolean equals(Object object) {
        if (!(object instanceof Randomizer that)) return false;
        return Objects.equals(version, that.version);
    }

    @Override
    public int hashCode() {
        return Objects.hashCode(version);
    }

    @Override
    public String toString() {
        return version.toString();
    }
}
