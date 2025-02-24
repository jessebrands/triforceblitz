package com.triforceblitz.triforceblitz.randomizer;

/**
 * Factory for creating generator instances.
 *
 * @see Generator
 */
public interface GeneratorFactory {
    Generator create(Randomizer randomizer, Preset preset);
}
