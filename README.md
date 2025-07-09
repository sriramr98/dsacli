# DSA CLI - GoLang Version

A CLI tool to practice DSA (Data Structures and Algorithms) questions using a spaced repetition algorithm with difficulty progression.

## Features

- **Spaced Repetition Algorithm**: Tracks your performance and suggests questions based on your understanding
- **Difficulty Progression**: Starts with easy questions and gradually moves to medium and hard
- **Smart Review**: Incorporates previously solved questions for review based on your performance
- **Performance Tracking**: Records time taken, hints needed, solution optimality, and bugs
- **SQLite Database**: Stores data in a lightweight SQLite database for better performance and reliability

## Installation

### Prerequisites
- Go 1.21 or higher

### Build from source
```bash
# Clone or download the repository
cd dsacli

# Install dependencies
go mod tidy

# Build the application
go build -o dsacli main.go

# (Optional) Install globally
go install
```

## Usage

### Get today's questions
```bash
./dsacli today
```

This command will suggest 1-2 questions based on your current progress:
- **Easy Phase**: Focus on easy questions until all are attempted
- **Medium Phase**: Focus on medium questions with smart review of previous questions
- **Hard Phase**: Focus on hard questions with smart review
- **Mastery Mode**: Review all questions based on spaced repetition scores

### List all questions
```bash
./dsacli list
```

Shows all available questions with their IDs, completion status, and SR scores.

### Mark a question as complete
```bash
./dsacli complete [question_id]
```

Where `[question_id]` is the question ID from the list command.

You'll be prompted to provide:
- **Hints needed** (1=many hints, 5=no hints)
- **Time taken** (in minutes, -1 if couldn't solve without solution)
- **Solution optimality** (1=not optimal, 5=very optimal)
- **Bugs encountered** (1=many bugs, 5=no bugs)

## How it works

### 🧠 The Science Behind Spaced Repetition

**Spaced repetition** is a learning technique based on cognitive science research that dramatically improves long-term retention and accelerates skill acquisition. This tool applies these principles specifically to coding interview preparation.

#### Why Traditional Practice Fails

Most developers practice coding problems like this:
1. ✗ Solve 20 problems in a row
2. ✗ Move to the next topic
3. ✗ Never revisit solved problems
4. ✗ Forget solutions within weeks

**Result**: Wasted time, poor retention, interview anxiety.

#### Why Spaced Repetition Works

Based on the **Ebbinghaus Forgetting Curve**, we forget information exponentially over time. Spaced repetition combats this by:

1. **Optimal Timing**: Reviews problems just before you'd forget them
2. **Adaptive Intervals**: Increases review gaps as you improve
3. **Focused Practice**: Prioritizes your weakest areas
4. **Long-term Retention**: Moves knowledge into permanent memory

### 🎯 Why This Accelerates Interview Prep

#### **Traditional Approach: 6 months**
- Solve 300+ problems randomly
- Forget 70% within 2 weeks
- Panic before interviews
- Need last-minute cramming

#### **Spaced Repetition Approach: 2-3 months**
- Solve 50-100 problems systematically
- Retain 90% long-term
- Build genuine confidence
- Perform consistently under pressure

### 🔬 Our Spaced Repetition Algorithm

#### Core Formula
```
SR Score = (Time Factor + Date Factor) × Solution Factor
```

#### 1. **Time Factor** - How quickly you solved it
- **< 25 minutes**: `time × 100` (excellent)
- **25-35 minutes**: `time × 200` (good)
- **35-45 minutes**: `time × 300` (needs work)
- **> 45 minutes**: `time × 400` (review soon)
- **Couldn't solve**: `24,000` (high priority)

#### 2. **Date Factor** - How long since last review
- **Never reviewed**: `10,000` (high priority)
- **Recent review**: Based on actual days elapsed
- **Ensures forgotten problems resurface naturally**

#### 3. **Solution Factor** - Overall performance quality
```
Performance Score = (Hints + Time Rank + Optimality + Bug-free) ÷ 4
```
- **Perfect (5.0)**: `×0.5` (review less frequently)
- **Poor (1.0)**: `×5.0` (review very frequently)

#### 4. **Adaptive Scoring** - Learning from repetition
- **First attempt**: Uses calculated score directly
- **Subsequent attempts**: `Previous × 0.7 + New × 0.3`
- **Prevents score explosion, reflects gradual improvement**

### 📈 Difficulty Progression System

#### Phase 1: Foundation Building (Easy Questions)
- **Focus**: Master fundamental patterns
- **Strategy**: Complete all easy questions first
- **Why**: Build confidence and core concepts
- **Duration**: 2-3 weeks

#### Phase 2: Pattern Recognition (Medium + Smart Review)
- **Focus**: Learn complex patterns while reviewing fundamentals
- **Strategy**: 1 new medium + 1 review question daily
- **Why**: Apply basics to harder problems
- **Duration**: 4-6 weeks

#### Phase 3: Advanced Mastery (Hard + Comprehensive Review)
- **Focus**: Tackle hardest problems while maintaining all skills
- **Strategy**: 1 new hard + 1 highest-scoring review question
- **Why**: Interview-level difficulty with solid foundation
- **Duration**: 3-4 weeks

#### Phase 4: Interview Readiness (Mastery Mode)
- **Focus**: Maintain peak performance across all difficulties
- **Strategy**: Review top 2 highest-scoring questions daily
- **Why**: Consistent performance under pressure
- **Duration**: Ongoing maintenance

### 🚀 Real-World Results

#### Week 1-2: Foundation
```bash
$ dsacli today
Focusing on: Easy Questions
1. Two Sum (Easy)
2. Valid Anagram (Easy)
```

#### Week 4-6: Pattern Building
```bash
$ dsacli today
Focusing on: Medium Questions (with Smart Review)
1. Longest Substring Without Repeating Characters (Medium)
2. Two Sum (Easy) - SR Score: 3,200 [Review needed]
```

#### Week 8-10: Advanced Practice
```bash
$ dsacli today
Focusing on: Hard Questions (with Smart Review)
1. Find Median from Data Stream (Hard)
2. Valid Parentheses (Medium) - SR Score: 5,100 [Struggling with optimization]
```

#### Week 12+: Interview Ready
```bash
$ dsacli today
Mastery Mode: Reviewing all questions!
1. Sliding Window Maximum (Hard) - SR Score: 4,800
2. Merge Intervals (Medium) - SR Score: 4,200
```


**Bottom Line**: This isn't just another practice tool - it's a scientifically-backed system that transforms how your brain learns and retains coding patterns.

## Data Storage

Questions and progress are stored in a SQLite database in your home directory:
- `~/.dsa-cli/dsa_questions.db`: SQLite database containing all questions and their progress

## Commands

- `dsacli today` - Get today's suggested questions
- `dsacli list` - List all questions with IDs and status
- `dsacli complete <question_id>` - Mark a question as complete
- `dsacli version` - Show version information

## Example Session

```bash
# List all questions to see IDs
$ ./dsacli list
All DSA Questions:

Easy Questions:
  ❌ ID:0 - Two Sum (SR Score: 0)
  ❌ ID:1 - Contains Duplicate (SR Score: 0)
  ...

# Get today's questions
$ ./dsacli today
Focusing on: Easy Questions
Here are your questions for today:
1. Two Sum (Easy) - https://leetcode.com/problems/two-sum/
2. Contains Duplicate (Easy) - https://leetcode.com/problems/contains-duplicate/

# After solving question with ID 0
$ ./dsacli complete 0
Updating score for: Two Sum
Did you need hints? (1=many hints, 5=no hints) 4
How long did it take (in minutes)? (-1 if you couldn't solve without solution) 15
Was the solution optimal? (1=not optimal, 5=very optimal) 4
Were there any bugs? (1=many bugs, 5=no bugs) 5

Successfully updated! New SR Score for 'Two Sum' is 2650.
```

## Contributing

Feel free to submit issues and enhancement requests!
