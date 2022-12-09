fn main() {
    println!("Hello, world!");
}




mod treetops {

    // Trees are visible from outside the grid if
    // they are the highest tree from their current
    // position to the edge in any of the four directions

    // While parsing, keep track of each of these values
    // for every point
    // hmax left
    // hmax right
    // vmax down
    // vmax up

    // Use async futures for values that don't exist yet
    // Then resolve everything
    // Then iterate through and count

    // Visibility from the left should always be known since we're parsing L->R
    // Visibility from the top should always be known since we're parsing TvB
    // Visibility from the right should be known by the time we finish the row
    

    // Or, visualize the numbers as stacked boolean layers
    // and use structured morphology operators?

    // Or use a convolution of some kind?

    // Or use pubsub with each tree listening for stuff it's interested in?

    // Or Map every point into 2 collections?
    // 


    // ==== Watershed
    // Start at each edge point
    // Drill inwards:
    //      
    //      
}