package main

import (
	"fmt"
	"time"
	"os"
	"bufio"
	"strings"
	"strconv"
)

const (
	layout = "2006-01-02T15:04:05"
)

type Play struct {
	name string 
	purchased []Ticket
	showStart time.Time
	showEnd time.Time 
}

type Comedy struct {
	laughs float32
	deaths int32
	Play
}

type Tragedy struct {
	laughs float32
	deaths int32
	Play
}

type Seat struct {
	number int32
	row int32
	cat *Category
}

type Category struct {
	name string
	base float32
}

type Ticket struct {
	customer string
	s *Seat
	show *Show
}

type Theatre struct {
	seats []Seat
	shows []Show
}

type Show interface {
	getName() string
	getShowStart() time.Time
	getShowEnd() time.Time
	addPurchase(*Ticket) bool
	isNotPurchased(*Ticket) bool
}

func (comedy *Comedy) getName() string {
	return comedy.name
}

func (tragedy *Tragedy) getName() string {
	return tragedy.name
}

func (comedy *Comedy) getShowStart() time.Time {
	return comedy.showStart
}

func (tragedy *Tragedy) getShowStart() time.Time {
	return tragedy.showStart
}

func (comedy *Comedy) getShowEnd() time.Time {
	return comedy.showEnd
}

func (tragedy *Tragedy) getShowEnd() time.Time {
	return tragedy.showEnd
}

func (comedy *Comedy) addPurchase(ticket *Ticket) bool {
	if comedy.isNotPurchased(ticket) {
		// A ticket for the same seat is NOT already sold. So we can the new ticket.
		comedy.purchased = append(comedy.purchased, *ticket)
		return true

	} else {
		return false 
	}
}

func (tragedy *Tragedy) addPurchase(ticket *Ticket) bool {
	if tragedy.isNotPurchased(ticket) {
		// A ticket for the same seat is NOT already sold. So we can the new ticket.
		tragedy.purchased = append(tragedy.purchased, *ticket)
		return true

	} else {
		return false 
	}
}

func (comedy *Comedy) isNotPurchased(ticket *Ticket) bool {
	seat_num := ticket.s.number
	comedy_tickets := comedy.purchased

	// Check if any of the tickets have the same seat_num
	for i := 0; i < len(comedy_tickets); i++ {
		if comedy_tickets[i].s.number == seat_num {
			// Seat already taken 
			return false
		}
	}

	// We iterated through the entire list and validated the ticket is for a new seat number
	return true
}

func (tragedy *Tragedy) isNotPurchased(ticket *Ticket) bool {
	seat_num := ticket.s.number
	tragedy_tickets := tragedy.purchased

	// Check if any of the tickets have the same seat_num
	for i := 0; i < len(tragedy_tickets); i++ {
		if tragedy_tickets[i].s.number == seat_num {
			// Seat already taken 
			return false
		}
	}

	// We iterated through the entire list and validated the ticket is for a new seat number
	return true
}

func NewComedy() Comedy {
	tickets := []Ticket{}
	dflt_start := "2020-03-03T16:00:00"
	dflt_end := "2020-03-03T17:20:00"
	show_start, _ := time.Parse(layout, dflt_start)
	show_end, _ := time.Parse(layout, dflt_end)
	comedy := Comedy{Play: Play{name: "Tartuffe", purchased: tickets, showStart: show_start, showEnd: show_end}}
	comedy.laughs = 0.2
	comedy.deaths = 0

	return comedy 
}

func NewTragedy() Tragedy{
	dflt_start := "2020-04-16T09:30:00"
	dflt_end := "2020-04-16T12:30:00"
	show_start, _ := time.Parse(layout, dflt_start)
	show_end, _ := time.Parse(layout, dflt_end)
	tickets := []Ticket{}
	tragedy := Tragedy{Play: Play{name: "Macbeth", purchased: tickets, showStart: show_start, showEnd: show_end}}
	tragedy.laughs = 0.0
	tragedy.deaths = 12

	return tragedy
}

func NewCategory() Category {
	category := Category{}
	category.name = "Standard"
	category.base = 25.0

	return category
}



// ****** Not sure if this is a solid practice or not. Essentially I am creating a separate method
// ****** for creating a category (other than NewCategory) where you actually pass values in 
// ****** and it then checks to see if they are valid values.
func CreateCategory(name string, base float32) Category {
	category := Category{}
	switch name {
		case "Prime":
			category.name = "Prime"

		case "Standard":
			category.name = "Standard"

		case "Special":
			category.name = "Special"

		default:
			fmt.Println("Invalid argument passed for category name - set to default value of \"standard\"")
			category.name = "Standard"
	}

	category.base = base

	return category 
}

func CreateTheatre(seats []Seat, shows []Show) Theatre{
	theatre := Theatre{}
	theatre.seats = seats
	theatre.shows = shows

	return theatre
}

func (tragedy *Tragedy) findRowFromSeatNum(seat_num int32) int32 {
	result := seat_num % 5
	if result == 0 {
		return 5
	} else {
		return result
	}
}

func (tragedy *Tragedy) checkSpecialTickets(seat_num int32) int32{
	for i := 5; i <= 25; i++ {
		if i % 5 == 0 && i != int (seat_num) {
			// Possible seat number in Special row
			tick := Ticket{}
			seat := Seat{}
			seat.number = int32 (i)
			tick.s = &seat
			if tragedy.isNotPurchased(&tick) {
				return int32 (i)
			} 
		}
	}

	// Could not find seat in Prime row
	return int32 (-1)
}

func (comedy *Comedy) checkSpecialTickets(seat_num int32) int32 {
	for i := 5; i <= 25; i++ {
		if i % 5 == 0 && i != int (seat_num) {
			// Possible seat number in Special row
			tick := Ticket{}
			seat := Seat{}
			seat.number = int32 (i)
			tick.s = &seat
			if comedy.isNotPurchased(&tick) {
				return int32 (i)
			} 
		}
	}

	// Could not find seat in Prime row
	return int32 (-1)
}

func (tragedy *Tragedy) checkStandardTickets(seat_num int32) int32 {
	for i := 2; i <= 25; i++ {
		mod := i % 5
		if mod == 2 || mod == 3 || mod == 4 && i != int (seat_num) {
			// Possible seat in Standard Row
			tick := Ticket{}
			seat := Seat{}
			seat.number = int32 (i)
			tick.s = &seat
			if tragedy.isNotPurchased(&tick) {
				return int32 (i)
			} 
		}
	}

	return int32 (-1)
}

func (comedy *Comedy) checkStandardTickets(seat_num int32) int32 {
	for i := 2; i <= 25; i++ {
		mod := i % 5
		if mod == 2 || mod == 3 || mod == 4 && i != int (seat_num) {
			// Possible seat in Standard Row
			tick := Ticket{}
			seat := Seat{}
			seat.number = int32 (i)
			tick.s = &seat
			if comedy.isNotPurchased(&tick) {
				return int32 (i)
			} 
		}
	}

	return int32 (-1)
}

func (tragedy *Tragedy) checkPrimeTickets(seat_num int32) int32 {
	for i := 1; i <= 21; i++ {
		if i % 5 == 1 && i != int (seat_num) {
			// Possible seat in Prime Row
			tick := Ticket{}
			seat := Seat{}
			seat.number = int32 (i)
			tick.s = &seat
			if tragedy.isNotPurchased(&tick) {
				return int32 (i)
			} 
		}
	}

	return int32 (-1)
}

func (comedy *Comedy) checkPrimeTickets(seat_num int32) int32 {
	for i := 1; i <= 21; i++ {
		if i % 5 == 1 && i != int (seat_num) {
			// Possible seat in Prime Row
			tick := Ticket{}
			seat := Seat{}
			seat.number = int32 (i)
			tick.s = &seat
			if comedy.isNotPurchased(&tick) {
				return int32 (i)
			} 
		}
	}

	return int32 (-1)
}

func (tragedy *Tragedy) findNewSeat(desired_seat_number int32) int32 {
	// Find an alternative seat in the same category
	// If can't, find more expensive.
	// If can't, find less expensive.
	if desired_seat_number % 5 == 0 {
		// Special Row
		special_seat := tragedy.checkSpecialTickets(desired_seat_number) 
		if  special_seat== -1 {
			prime_seat := tragedy.checkPrimeTickets(desired_seat_number)
			if prime_seat == -1 {
				return tragedy.checkStandardTickets(desired_seat_number)
			} else {
				return prime_seat
			}
		} else {
			return special_seat
		}
	} else if desired_seat_number % 5 == 1 {
		// Prime Row
		prime_seat := tragedy.checkPrimeTickets(desired_seat_number)
		if prime_seat == -1 {
			std_seat := tragedy.checkStandardTickets(desired_seat_number)
			if std_seat == -1 {
				return tragedy.checkSpecialTickets(desired_seat_number)
			} else {
				return std_seat
			}
		} else {
			return prime_seat
		}
	} else {
		// Standard Row
		std_seat := tragedy.checkStandardTickets(desired_seat_number)
		if std_seat == -1 {
			prime_seat := tragedy.checkPrimeTickets(desired_seat_number)
			if prime_seat == -1 {
				return tragedy.checkSpecialTickets(desired_seat_number)
			} else {
				return prime_seat
			}
		} else {
			return std_seat
		}
	}
}

func (comedy *Comedy) findNewSeat(desired_seat_number int32) int32 {
	if desired_seat_number % 5 == 0 {
		// Special Row
		special_seat := comedy.checkSpecialTickets(desired_seat_number) 
		if  special_seat== -1 {
			prime_seat := comedy.checkPrimeTickets(desired_seat_number)
			if prime_seat == -1 {
				return comedy.checkStandardTickets(desired_seat_number)
			} else {
				return prime_seat
			}
		} else {
			return special_seat
		}
	} else if desired_seat_number % 5 == 1 {
		// Prime Row
		prime_seat := comedy.checkPrimeTickets(desired_seat_number)
		if prime_seat == -1 {
			std_seat := comedy.checkStandardTickets(desired_seat_number)
			if std_seat == -1 {
				return comedy.checkSpecialTickets(desired_seat_number)
			} else {
				return std_seat
			}
		} else {
			return prime_seat
		}
	} else {
		// Standard Row
		std_seat := comedy.checkStandardTickets(desired_seat_number)
		if std_seat == -1 {
			prime_seat := comedy.checkPrimeTickets(desired_seat_number)
			if prime_seat == -1 {
				return comedy.checkSpecialTickets(desired_seat_number)
			} else {
				return prime_seat
			}
		} else {
			return std_seat
		}
	}
} 

// ****** Not sure what he means in the assignment for this function. He states seat and row
// ****** have default values of 1. Then how can we enforce this if we are getting parameters passed
// ****** into the function? 
func NewSeat(seat_number int32, row_number int32, category *Category) Seat {
	seat := Seat{}
	seat.number = seat_number
	seat.row = row_number
	seat.cat = category

	return seat
}

func NewTicket(cust_name string, seat *Seat, show *Show) Ticket {
	ticket := Ticket{}
	ticket.customer = cust_name
	ticket.s = seat
	ticket.show = show

	return ticket
}

func NewTheatre(num_seats int32, shows []Show) Theatre {
	theatre := Theatre{}
	seats := make([]Seat, num_seats)
	theatre.seats = seats
	theatre.shows = shows

	return theatre
}

func main() {
	// Create the theatre.
	theatre := Theatre{}
	var row, col, seat_count int32
	seat_count = 1
	seats := make([]Seat, 25)

	// Create the seats
	for col = 1; col <= 5; col++ {
		for row = 1; row <= 5; row++ {
			if row == 1 {
				category := CreateCategory("Prime", 35.0)
				seat := NewSeat(seat_count, row, &category)
				seats[seat_count - 1] = seat
				seat_count++
				
			} else if (row == 2 || row == 3 || row == 4) {
				category := CreateCategory("Standard", 25.0)
				seat := NewSeat(seat_count, row, &category)
				seats[seat_count - 1] = seat
				seat_count++

			} else {
				category := CreateCategory("Special", 15.0)
				seat := NewSeat(seat_count, row, &category)
				seats[seat_count - 1] = seat
				seat_count++
			}
		}
	}

	// Add all the seats to the theatre.
	theatre.seats = seats

	// Create the shows
	shows := make([]Show, 0)

	// Comedy Show
	comedy := NewComedy()
	comedy_start := "2020-03-03T19:30:00"
	comedy_end := "2020-03-03T20:00:00"
	comedy.showStart, _ = time.Parse(layout, comedy_start)
	comedy.showEnd, _ = time.Parse(layout, comedy_end)

	// Tragedy Show
	tragedy := NewTragedy()
	tragedy_start := "2020-04-10T20:00:00"
	tragedy_end := "2020-04-10T23:00:00"
	tragedy.showStart, _ = time.Parse(layout, tragedy_start)
	tragedy.showEnd, _ = time.Parse(layout, tragedy_end)

	// Add the shows to the show list	
	shows = append(shows, &comedy)
	shows = append(shows, &tragedy)

	// Add the shows to the theatre
	theatre.shows = shows


	// ** Console application **

	fmt.Println("\n", "********** Welcome to the CSI 2120 Theatre! **********", "\n")

	show_names := theatre.shows

	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Please enter the play you would like to see: ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		input = strings.TrimRight(input, "\r\n")
		count := 0
		for i := 0; i < len(show_names); i++ {
			if show_names[i].getName() == input {
				fmt.Println("\n", "Please enter a seat number for the", input + " show:")
				ticket_input, _ := reader.ReadString('\n')
				ticket_input = strings.TrimRight(ticket_input, "\r\n")
				ticket_input = strings.TrimSpace(ticket_input)
				seat_num_input, err := strconv.Atoi(ticket_input)
				if err == nil && seat_num_input <= 25 {
					seat_number := int32 (seat_num_input)
					seat := Seat{}
					seat.number = seat_number
					ticket := Ticket{}
					ticket.s = &seat

					// Two cases for the separate shows
					if strings.TrimRight(input, "\r\n") == "Macbeth" {
						// Tragedy 
						if tragedy.isNotPurchased(&ticket) {
							tragedy.addPurchase(&ticket)
							fmt.Println("Successfully purchased seat number: ", seat_number)
						} else {
							// Recommend another seat
							new_seat := tragedy.findNewSeat(seat_number)
							if new_seat == -1 {
								fmt.Println("Unfortunately there are no seats available!")
							} else {
								fmt.Println("We found seat number ", new_seat, " as an alternative!")
							}
						}

					} else if strings.TrimRight(input, "\r\n") == "Tartuffe" {
						// Comedy 
						if comedy.isNotPurchased(&ticket) {
							comedy.addPurchase(&ticket)
							fmt.Println("Successfully purchased seat number: ", seat_number)
						} else {
							// Recommend another seat
							new_seat := comedy.findNewSeat(seat_number)
							if new_seat == -1 {
								fmt.Println("Unfortunately there are no seats available!")
							} else {
								fmt.Println("We found seat number ", new_seat, " as an alternative!")
							}
						}
					}
				} else {
					fmt.Println("Unfortunately ", ticket_input, " is not a value we accept for seat numbers.")
				}
			} else {
				count++
			}
		}
		if count == len(show_names) {
			fmt.Println("Unfortunately the CSI 2120 Theatre isn't offering that play right now. Try again with a different name.", "\n")
			count = 0
		}
	}
}