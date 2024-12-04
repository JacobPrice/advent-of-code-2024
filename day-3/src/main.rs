use anyhow::Result;
use dotenv::dotenv;
use std::env;
use regex::Regex;
use std::ops::Range;

fn find_next_dont_index(do_idx: &usize, dont_indices: &[usize]) -> Option<usize> {
    dont_indices.iter()
        .find(|&&idx| idx > *do_idx)
        .copied()
}

fn main() -> Result<()> {
    // Load environment variables from .env file
    dotenv().ok();

    // Get environment variables
    let cookie = env::var("AOC_COOKIE").expect("AOC_COOKIE must be set in .env file");
    let url = env::var("AOC_URL").expect("AOC_URL must be set in .env file");

    // Create a client and set the cookie
    let client = reqwest::blocking::Client::new();
    let input = client
        .get(url)
        .header("Cookie", format!("session={}", cookie))
        .send()?
        .text()?;

    // Create regex patterns
    let mul_re = Regex::new(r"mul\((\d+),(\d+)\)")?;
    let do_re = Regex::new(r"do\(\)")?;
    let dont_re = Regex::new(r"don't\(\)")?;

    // Collect all do() positions (including 0 for initial state)
    let mut do_indices = vec![0];
    for capture in do_re.captures_iter(&input) {
        do_indices.push(capture.get(0).unwrap().start());
    }

    // Collect all don't() positions
    let mut dont_indices = vec![];
    for capture in dont_re.captures_iter(&input) {
        dont_indices.push(capture.get(0).unwrap().start());
    }

    // Create ranges where multiplications are enabled
    let input_len = input.len();
    let mut enabled_ranges = vec![];
    for do_idx in do_indices {
        let next_dont = find_next_dont_index(&do_idx, &dont_indices).unwrap_or(input_len);
        enabled_ranges.push(do_idx..next_dont);
    }

    // Process multiplications
    let mut sum = 0;
    for capture in mul_re.captures_iter(&input) {
        let pos = capture.get(0).unwrap().start();
        let num1: u32 = capture[1].parse().unwrap();
        let num2: u32 = capture[2].parse().unwrap();

        // Check if this multiplication falls within any enabled range
        if enabled_ranges.iter().any(|range| range.contains(&pos)) {
            let result = num1 * num2;
            println!("Processing enabled multiplication at {}: {}*{} = {}", pos, num1, num2, result);
            sum += result;
        } else {
            println!("Skipping disabled multiplication at {}: {}*{}", pos, num1, num2);
        }
    }

    println!("Final sum of all enabled multiplications: {}", sum);

    Ok(())
}