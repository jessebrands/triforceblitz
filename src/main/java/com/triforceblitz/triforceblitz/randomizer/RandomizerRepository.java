package com.triforceblitz.triforceblitz.randomizer;

import java.util.List;
import java.util.Optional;

/**
 * Repository for loading randomizers.
 *
 * @see Randomizer
 */
public interface RandomizerRepository {
    /**
     * Loads a single randomizer.
     *
     * @param version Version of the randomizer to load.
     * @return {@link Optional} containing the loaded randomizer, or empty.
     */
    Optional<Randomizer> load(RandomizerVersion version);

    /**
     * Loads all Randomizers off the
     *
     * @return {@link List} of loaded randomizers.
     */
    List<Randomizer> loadAll();

    /**
     * Checks if a randomizer exists in the repository.
     *
     * @param version Version of the randomizer.
     * @return <code>true</code> if it exists, <code>false</code> if not.
     */
    boolean exists(RandomizerVersion version);
}
