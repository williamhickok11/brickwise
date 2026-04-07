export interface ParsedREPSEntry {
  hours: number
  category: string
  description: string
  date: string // YYYY-MM-DD format
  mileage: number
  full_drive: boolean
  confidence: 'high' | 'medium' | 'low'
}

/**
 * Parses free-form text input to extract REPS time entry information.
 * Handles various formats like:
 * - "2.5 hours of property management fixing the sink"
 * - "Spent 3 hours on maintenance and repairs, drove 22 miles"
 * - "Today: 1.5h contractor oversight at Oak Avenue"
 * - "Yesterday: 2 hours accounting and admin"
 */
export function parseREPSText(text: string, defaultDate: string = getTodayDate()): ParsedREPSEntry {
  const normalized = text.toLowerCase().trim()
  
  // Default values
  let hours = 0
  let category = ''
  let date = defaultDate
  let mileage = 0
  let full_drive = false
  let confidence: 'high' | 'medium' | 'low' = 'low'
  let hasDateCue = false
  let hasHours = false

  // Extract date references
  const dateMatch = normalized.match(/\b(today|yesterday|tomorrow)\b/)
  if (dateMatch) {
    hasDateCue = true
    const today = new Date()
    if (dateMatch[0] === 'yesterday') {
      today.setDate(today.getDate() - 1)
      date = formatDate(today)
    } else if (dateMatch[0] === 'tomorrow') {
      today.setDate(today.getDate() + 1)
      date = formatDate(today)
    } else {
      date = defaultDate
    }
  }

  // Extract hours - look for patterns like "2.5 hours", "2h", "2 hrs", "2.5h"
  const hourPatterns = [
    /(\d+\.?\d*)\s*(?:hours?|hrs?|h)\b/i,
    /\b(\d+\.?\d*)\s*h\b/i,
  ]
  
  for (const pattern of hourPatterns) {
    const match = normalized.match(pattern)
    if (match) {
      hours = parseFloat(match[1])
      if (hours > 0 && hours <= 24) {
        hasHours = true
        break
      }
    }
  }

  // Extract mileage - look for patterns like "22 miles", "22 mi", "drove 22"
  const mileagePatterns = [
    /(\d+\.?\d*)\s*(?:miles?|mi)\b/i,
    /drove\s+(\d+\.?\d*)/i,
    /mileage[:\s]+(\d+\.?\d*)/i,
  ]

  for (const pattern of mileagePatterns) {
    const match = normalized.match(pattern)
    if (match) {
      mileage = parseFloat(match[1])
      if (mileage > 0) {
        full_drive = true
        break
      }
    }
  }

  // Confidence semantics prioritize whether we captured valid hours.
  if (hasHours) {
    confidence = hasDateCue ? 'high' : 'medium'
  } else if (hasDateCue) {
    confidence = 'medium'
  }

  // Clean up description - remove extracted info to make it cleaner
  let cleanDescription = text.trim()
  
  // Remove date references
  cleanDescription = cleanDescription.replace(/\b(today|yesterday|tomorrow)[:\s]*/gi, '')
  
  // Remove hour patterns
  cleanDescription = cleanDescription.replace(/(\d+\.?\d*)\s*(?:hours?|hrs?|h)\b/gi, '')
  
  // Remove mileage patterns
  cleanDescription = cleanDescription.replace(/(\d+\.?\d*)\s*(?:miles?|mi)\b/gi, '')
  cleanDescription = cleanDescription.replace(/drove\s+(\d+\.?\d*)/gi, '')
  
  // Clean up extra whitespace
  cleanDescription = cleanDescription.replace(/\s+/g, ' ').trim()
  
  // If description is empty or too short, use original text
  if (cleanDescription.length < 10) {
    cleanDescription = text.trim()
  }

  return {
    hours,
    category,
    description: cleanDescription || text.trim(),
    date,
    mileage,
    full_drive,
    confidence,
  }
}

function getTodayDate(): string {
  return formatDate(new Date())
}

function formatDate(date: Date): string {
  return date.toISOString().split('T')[0]
}
