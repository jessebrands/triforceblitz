package com.triforceblitz.triforceblitz.randomizer;

import java.util.regex.Pattern;

/**
 * Ocarina of Time Randomizer version value.
 *
 * @see Randomizer
 */
public record RandomizerVersion(
        int major,
        int minor,
        int patch,
        String branch,
        int branchMajor,
        int branchMinor
) {
    private static final Pattern VALID_VERSION = Pattern.compile(
            "^v?(0|[1-9][0-9]*)\\.(0|[1-9][0-9]*)\\.(0|[1-9][0-9]*)-([a-z][a-z0-9]*)-(0|[1-9][0-9]*)\\.(0|[1-9][0-9]*)$"
    );

    /**
     * Creates a new version value from a string.
     *
     * @param version Randomizer version string.
     * @return New instance of {@link RandomizerVersion}.
     */
    public static RandomizerVersion of(String version) {
        var matcher = VALID_VERSION.matcher(version);
        if (!matcher.matches()) {
            throw new IllegalArgumentException("not a valid Ocarina of Time Randomizer version");
        }
        return new RandomizerVersion(
                Integer.parseInt(matcher.group(1)),
                Integer.parseInt(matcher.group(2)),
                Integer.parseInt(matcher.group(3)),
                matcher.group(4),
                Integer.parseInt(matcher.group(5)),
                Integer.parseInt(matcher.group(6))
        );
    }

    /**
     * Checks if a string is a valid Randomizer version.
     *
     * @param version Version string.
     * @return <code>true</code> if valid, <code>false</code> if invalid.
     */
    public static boolean valid(String version) {
        return VALID_VERSION.matcher(version).matches();
    }

    @Override
    public String toString() {
        return String.format(
                "%d.%d.%d-%s-%d.%d",
                major,
                minor,
                patch,
                branch,
                branchMajor,
                branchMinor
        );
    }
}
